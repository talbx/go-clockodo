package intercept

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/talbx/go-clockodo/pkg/model"
)

type Interceptor interface {
	Intercept()
}

type ConfigReaderInterceptor struct{}

var ClockodoConfig model.GoClockodoConfig

func (i ConfigReaderInterceptor) Intercept() {
	ReadConfig(&ClockodoConfig)
}

func ReadConfig(config *model.GoClockodoConfig) {
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	config.ApiKey = viper.GetString("apiKey")
	config.ApiUser = viper.GetString("apiUser")
}
