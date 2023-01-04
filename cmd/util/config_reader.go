package util

import (
	"fmt"

	"github.com/spf13/viper"
)

type GoClockodoConfig struct {
	ApiKey string
	ApiUser string
}


func ReadConfig(config *GoClockodoConfig) {
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	config.ApiKey = viper.GetString("apiKey")
	config.ApiUser = viper.GetString("apiUser")
	fmt.Println(viper.Get("apiKey")) // this would be "steve"
	fmt.Println(viper.Get("apiUser")) // this would be "steve"
}
