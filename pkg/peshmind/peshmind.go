package peshmind

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Debug                 bool              `json:"debug"`
	EndPoint              string            `json:"endpoint"`
	Switches              map[string]Switch `json:"switches"`
	SimGeneratePercentage int               `json:"sim_generate_percentage"`
	Simulations           []string          `json:"simulations"`
	UserNameDefault       string            `json:"username_default"`
	PasswordDefault       string            `json:"password_default"`
}

func NewConfig() *Config {
	return &Config{
		Switches:              make(map[string]Switch),
		UserNameDefault:       "admin",
		PasswordDefault:       "ask",
		Simulations:           make([]string, 0),
		SimGeneratePercentage: 100,
	}

}

// Save the current configuration to a JSON file
func (c *Config) SaveConfig(path string) error {
	// Check if the path is already existing
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file already exists: %s", path)
	}

	configBytes, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	err = os.WriteFile(path, configBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config to file: %w", err)
	}
	return nil
}

// LoadConfig loads the configuration from a JSON file
func (c *Config) LoadConfig(path string) error {
	// Check if the path is not existing
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", path)
	}
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	err = json.Unmarshal(configBytes, c)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return nil
}

// ShowConfig prints the current configuration to the console
func (c *Config) ShowConfig() {
	fmt.Println("Current configuration:")
	for key, value := range c.Switches {
		fmt.Printf("%s: %s\n", key, value)
	}
}
