package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Env            string `json:"env"`
	SendGridAPIKey string `json:"send_grid_api_key"`
	EmailFrom      string `json:"email_from"`
}

func New() (*Config, error) {
	conf := &Config{}
	err := godotenv.Load("../.env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %s", err)
	}
	err = conf.loadConfigEnv()
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (conf *Config) loadConfigEnv() error {
	getEnvironment := func(data []string, getkeyVal func(item string) (key, val string)) map[string]string {
		items := make(map[string]string)
		for _, item := range data {
			key, val := getkeyVal(item)
			items[key] = val
		}
		return items
	}

	env := getEnvironment(os.Environ(), func(item string) (key, val string) {
		splits := strings.Split(item, "=")
		key = splits[0]
		val = splits[1]
		return
	})
	jsonEnv := make(map[string]string)
	for key, val := range env {
		jsonEnv[strings.ToLower(key)] = val

	}
	jsonData, _ := json.Marshal(jsonEnv)
	err := json.Unmarshal(jsonData, conf)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: %s", err.Error())
	}
	return nil
}
