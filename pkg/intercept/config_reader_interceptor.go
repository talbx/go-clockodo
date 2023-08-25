package intercept

import (
	"github.com/spf13/viper"
	"github.com/talbx/go-clockodo/pkg/model"
)

type Interceptor interface {
	Intercept() error
}

type ConfigReaderInterceptor struct{}

var ClockodoConfig model.GoClockodoConfig

func (i ConfigReaderInterceptor) Intercept() error {
	return ReadConfig(&ClockodoConfig)
}

func ReadConfig(config *model.GoClockodoConfig) error {
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		return err
	}
	config.ApiKey = viper.GetString("apiKey")
	config.ApiUser = viper.GetString("apiUser")
	return nil
}
