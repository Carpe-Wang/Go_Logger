package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

func main() {
	viper.SetDefault("filedir", "./")
	viper.SetConfigName("config")       //Configure file name and no extension
	viper.SetConfigType("yaml")         //Configure file type
	viper.AddConfigPath("/etc/appname") //Find configuration file path
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
		panic(fmt.Errorf("fatal error config File:", err))
	}

	//Real time monitoring of file changes
	viper.WatchConfig()
	//if happen changes
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
	})

}
