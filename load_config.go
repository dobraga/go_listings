package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	log.Debug("Loading .env")
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	log.Debug("Loaded .env")

	log.Debug("Loading settings.toml")
	viper.SetConfigFile("settings.toml")
	err = viper.MergeInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	log.Debug("Loaded settings.toml")

	env := viper.Get("ENV")

	settings_file := fmt.Sprintf("settings_%s.toml", env)
	log.Debug(fmt.Sprintf("Loading %s", settings_file))

	viper.SetConfigFile(settings_file)
	err = viper.MergeInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	log.Debug(fmt.Sprintf("Loaded %s", settings_file))

	SQLALCHEMY_DATABASE_URI := viper.Get("SQLALCHEMY_DATABASE_URI").(string)
	POSTGRES_USER := viper.Get("POSTGRES_USER")
	if POSTGRES_USER == nil {
		panic("Need 'POSTGRES_USER' variable in .env file")
	}

	POSTGRES_PASSWORD := viper.Get("POSTGRES_PASSWORD")
	if POSTGRES_PASSWORD == nil {
		panic("Need 'POSTGRES_PASSWORD' variable in .env file")
	}

	POSTGRES_DB := viper.Get("POSTGRES_DB")
	if POSTGRES_DB == nil {
		panic("Need 'POSTGRES_DB' variable in .env file")
	}

	SQLALCHEMY_DATABASE_URI = fmt.Sprintf(SQLALCHEMY_DATABASE_URI, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB)

	viper.Set("SQLALCHEMY_DATABASE_URI", SQLALCHEMY_DATABASE_URI)
}
