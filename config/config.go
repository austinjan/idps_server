package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

//VERSION application version
var VERSION = "0.0.1"

func setDefault() {
	viper.SetDefault("host", "localhost")
	viper.SetDefault("webserver", map[string]interface{}{
		"port":    ":3011",
		"timeout": 5,
	})
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.idps_server")
	viper.SetDefault("version", VERSION)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		fmt.Printf("Fatal error config file: %s use default config\n", err)
		setDefault()
		viper.WriteConfigAs("./config.json")
		fmt.Println("Writting default config file")
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed: ", e.Name)
		viper.ReadInConfig()
	})
}
