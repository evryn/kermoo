package user_config

import (
	"bufio"
	"buggybox/modules/logger"
	"buggybox/modules/utils"
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"
)

func MustLoadProvidedConfig(filename string) {
	uc, err := LoadUserConfig(filename)

	if err != nil {
		logger.Log.Fatal("invalid initial user-provided config", zap.Error(err))
	}

	logger.Log.Info("initial configuration is loaded", zap.Any("unprepared", *uc))

	prepared, err := uc.GetPreparedConfig()

	if err != nil {
		logger.Log.Fatal("invalid prepared user-provided config", zap.Error(err))
	}

	logger.Log.Info("prepared configuration is loaded", zap.Any("prepared", prepared))

	Prepared = *prepared
}

func LoadUserConfig(filename string) (*UserConfigType, error) {
	var err error
	var uc UserConfigType

	if filename == "" {
		return nil, fmt.Errorf("provided filename is empty")
	}

	if filename == "-" {
		logger.Log.Debug("loading configuration from stdin...")

		content, err := readStdin()
		if err != nil {
			return nil, err
		}

		uc, err = unmarshal(content)
		if err != nil {
			return nil, err
		}
	} else {
		logger.Log.Debug("loading configuration from file...", zap.String("filename", filename))
		content, err := readFile(filename)
		if err != nil {
			return nil, err
		}

		uc, err = unmarshal(content)
		if err != nil {
			return nil, err
		}
	}

	err = uc.Validate()

	if err != nil {
		return nil, err
	}

	return &uc, nil
}

func readFile(filename string) (string, error) {
	logger.Log.Debug("reading file", zap.String("filename", filename))

	body, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("unable to read file: %s", filename)
	}

	return string(body), nil
}

func readStdin() (string, error) {
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
