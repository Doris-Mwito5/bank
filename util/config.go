package util

import (
	"github.com/spf13/viper"
)

// Config stores all the configuration var
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

//read configuration from file or environment var
func LoadConfig(path string) (config Config, err error) {
	//tell viper the location of the config file
	viper.AddConfigPath(path)
	//tell viper to look for the specified name
	viper.SetConfigName("app")
	//tell viper to look for the specified type
	viper.SetConfigType("env")

	//tell viper to read the env values
	viper.AutomaticEnv()

	//start reading the config var
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	//unmarshall the env var to the config object
	err = viper.Unmarshal(&config)
	return
}
