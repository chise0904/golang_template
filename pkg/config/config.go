package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// LoadConfig read config
func LoadConfig(path string, cfg interface{}) error {
	// set file type toml or yaml
	viper.AutomaticEnv()
	// check user want setting other config
	name := viper.GetString("CONFIG_NAME")
	if name == "" {
		name = "config"
	}

	viper.SetConfigType("toml")
	viper.SetConfigName(name)
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Error().Msgf("Error reading config file, %s", err.Error())

		return err
	}

	err := viper.Unmarshal(cfg)
	if err != nil {
		log.Error().Msgf("unable to decode into struct, %v", err.Error())
		return err
	}

	return nil
}
