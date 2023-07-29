package user_config

import (
	"bufio"
	"buggybox/modules/logger"
	"buggybox/modules/planner"
	"buggybox/modules/process"
	"buggybox/modules/utils"
	"buggybox/modules/web_server"
	"encoding/json"
	"log"
	"os"

	"go.uber.org/zap"
)

var UserConfig UserConfigType

type UserConfigType struct {
	ApiVersion string                 `json:"apiVersion"`
	Process    process.Process        `json:"process"`
	Plans      []planner.Plan         `json:"plans"`
	WebServers []web_server.WebServer `json:"webServers"`
}

func MustLoadUserConfig(filename string) {
	if filename == "" {
		logger.Log.Fatal("provided filename is empty")
	}

	if filename == "-" {
		logger.Log.Debug("loading configuration from stdin...")
		UserConfig = mustUnmarshal(mustReadStdin())
		return
	}

	logger.Log.Debug("loading configuration from file...", zap.String("filename", filename))
	UserConfig = mustUnmarshal(mustReadFile(filename))
}

func mustReadFile(filename string) string {
	logger.Log.Debug("reading file", zap.String("filename", filename))

	body, err := os.ReadFile(filename)
	if err != nil {
		logger.Log.Fatal("unable to read file.", zap.String("filename", filename))
	}

	return string(body)
}

func mustReadStdin() string {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		logger.Log.Fatal("stdin is not available to read from")
	}

	var stdin []byte
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		stdin = append(stdin, scanner.Bytes()...)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return string(stdin)
}

func mustUnmarshal(content string) UserConfigType {

	firstChar := string(content[0])

	// Content is probably JSON
	if firstChar == "{" {
		return mustUnmarshalJson(content)
	}

	// Content is probably YAML otherwise. Convert to json and unmarshal
	logger.Log.Debug("converting yaml to json", zap.String("yaml", content))

	json, err := utils.YamlToJSON(content)

	if err != nil {
		logger.Log.Fatal("invalid yaml configuration. check syntax errors.", zap.Error(err))
	}

	return mustUnmarshalJson(json)
}

func mustUnmarshalJson(jsonContent string) UserConfigType {
	uc := UserConfig

	logger.Log.Debug("unmarshalling json content", zap.String("json", jsonContent))

	err := json.Unmarshal([]byte(jsonContent), &uc)
	if err != nil {
		logger.Log.Fatal("unable to unmarshal json content", zap.Error(err))
	}

	return uc
}
