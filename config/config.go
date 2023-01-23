package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfig() error {
	viper.SetConfigType("json")

	// load env config
	viper.SetConfigFile(ENV_CONFIG_FILE)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}
	if err := viper.Unmarshal(&appCfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// load db config
	viper.SetConfigFile(DB_CONFIG_FILE)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}
	if err := viper.Unmarshal(&dbCfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return nil
}
