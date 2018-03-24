package config

import (
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

func Load() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicf("Fatal error config file: %s \n", err)
	}
}