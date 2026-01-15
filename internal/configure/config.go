package configure

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	CONFIG_FILE_NAME string = "connection_config.json"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUsername string `json:"current_username"`
}

func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error home directory: %w", err)
	}
	filePath := filepath.Join(homeDir, CONFIG_FILE_NAME)
	return filePath, nil
}

func Read() (Config, error) {
	filePath, err := GetConfigPath()
	if err != nil {
		return Config{}, fmt.Errorf("error getting config file path: %w", err)
	}
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("error reading file from path=%v :%w", filePath, err)
	}
	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshalling bytes into config struct: %w", err)
	}
	return config, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUsername = username
	filePath, err := GetConfigPath()
	if err != nil {
		fmt.Errorf("error getting config file path: %w", err)
	}
	bytes, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("error marshalling object=%v: %w", c, err)
	}
	err = os.WriteFile(filePath, bytes, 0644) // 0644 is READ-WRITE permission
	if err != nil {
		return fmt.Errorf("error writing bytes to file path=%v: %w", filePath, err)
	}
	return nil
}
