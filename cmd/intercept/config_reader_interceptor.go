package intercept

import (
	"fmt"

	"github.com/spf13/viper"
)

type Interceptor interface {
	Intercept()
}

type ConfigReaderInterceptor struct{}

var ClockodoConfig GoClockodoConfig

func (i ConfigReaderInterceptor) Intercept() {
	ReadConfig(&ClockodoConfig)
}

type GoClockodoConfig struct {
	ApiKey string
	ApiUser string
	TestCustomer int
	MainCustomer int
	TestService int
	TestProject int
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
	config.TestCustomer = viper.GetInt("testCustomerId")
	config.TestService = viper.GetInt("testServiceId")
	config.TestProject = viper.GetInt("testProjectId")
}
