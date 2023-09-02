package user_config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"kermoo/modules/logger"
	"kermoo/modules/utils"
	"os"
	"os/user"
	"strings"

	"go.uber.org/zap"
)

func MustLoadPreparedConfig(config string) {
	prepared, err := MakePreparedConfig(config)

	if err != nil {
		logger.Log.Panic(err.Error())
	}

	Prepared = *prepared
}

// MakePreparedConfig resolves the configuration content and prepares the config object in the following order:
//
// 1. When no config is provided, it autoloads it from `{userHome}/.kermoo/config.[yaml|yml|json]`
//
// 2. When it couldn't autoload from file, it tries to load from `KERMOO_CONFIG` environment variable
//
// 3. When config is given and equals to "-", it tries to read from stdin (Standard Input) pipe
//
// 4. When the config is a file path, it tries to load it from that file.
//
// 5. Otherwise, it considers the content of config as the actual config and tries to parse that.
func MakePreparedConfig(config string) (*PreparedConfigType, error) {
	uc, err := makeUserConfig(config)

	if err != nil {
		return nil, fmt.Errorf("invalid config: %v", err)
	}

	prepared, err := uc.GetPreparedConfig()

	if err != nil {
		return nil, fmt.Errorf("unable to preapre parsed config: %v", err)
	}

	return prepared, nil
}

func makeUserConfig(config string) (*UserConfigType, error) {
	var err error
	var uc UserConfigType

	config, err = getResolvedConfig(config)

	if err != nil {
		return nil, err
	}

	if config == "" {
		return nil, fmt.Errorf("resolved config is empty")
	}

	uc, err = unmarshal(config)
	if err != nil {
		return nil, err
	}

	err = uc.Validate()

	if err != nil {
		return nil, err
	}

	return &uc, nil
}

func getResolvedConfig(config string) (string, error) {
	if config == "" {
		return getAutoloadedConfig()
	}

	if config == "-" {
		logger.Log.Debug("loading configuration from stdin...")

		return readStdin()
	}

	// TODO: Change the following to something more sophisticated
	if !strings.ContainsAny(config, "\n\t{}") {
		return readFile(config)
	}

	return config, nil
}

func getAutoloadedConfig() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("unable to determine current user to autoload config: %v", err)
	}

	pathes := []string{
		u.HomeDir + "/.kermoo/config.yaml",
		u.HomeDir + "/.kermoo/config.yml",
		u.HomeDir + "/.kermoo/config.json",
		u.HomeDir + "/.config/kermoo/config.yaml",
		u.HomeDir + "/.config/kermoo/config.yml",
		u.HomeDir + "/.config/kermoo/config.json",
	}

	for _, path := range pathes {
		content, err := readFile(path)

		if err == nil {
			return content, nil
		}
	}

	content := os.Getenv("KERMOO_CONFIG")
	if content != "" {
		return content, nil
	}

	return "", fmt.Errorf("no config is specified so we tried to autoload config but was unable to load it either from default home paths (like %v) or from the environment variable. kermoo can not live without a config :(", pathes[0])
}

func readFile(filename string) (string, error) {
	logger.Log.Debug("reading file ...", zap.String("filename", filename))

	body, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("unable to read file: %s", filename)
	}

	return string(body), nil
}

func readStdin() (string, error) {
	logger.Log.Debug("reading stdin ...")

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", fmt.Errorf("stdin is not available to read from")
	}

	var stdin string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		stdin += string(line) + "\n"
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return string(stdin), nil
}

func unmarshal(content string) (UserConfigType, error) {
	logger.Log.Debug("unmarshaling config ...", zap.Any("config", content))

	firstChar := string(content[0])

	// Content is probably JSON
	if firstChar == "{" {
		return unmarshalJson(content)
	}

	// Content is probably YAML otherwise. Convert to json and unmarshal
	logger.Log.Debug("converting yaml to json", zap.String("yaml", content))

	json, err := utils.YamlToJSON(content)

	if err != nil {
		return UserConfigType{}, fmt.Errorf("invalid yaml configuration - check syntax errors. %w", err)
	}

	return unmarshalJson(json)
}

func unmarshalJson(jsonContent string) (UserConfigType, error) {
	uc := UserConfigType{}

	logger.Log.Debug("unmarshalling json content", zap.String("json", jsonContent))

	err := json.Unmarshal([]byte(jsonContent), &uc)
	if err != nil {
		return UserConfigType{}, fmt.Errorf("unable to unmarshal json content. %w", err)
	}

	return uc, nil
}
