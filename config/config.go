package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

)
//VERSION application version
var VERSION = "0.0.1"

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.idps_server")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		fmt.Printf("Fatal error config file: %s use default config\n", err)
		viper.WriteConfigAs("./eyesfreeconfig.json")
		fmt.Println("Writting default config file")
	}
	viper.SetDefault("version", VERSION)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed: ", e.Name)
		viper.ReadInConfig()
	})
}
