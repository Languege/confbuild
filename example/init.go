package example

import (
	"github.com/spf13/viper"
	"fmt"
)


func init() {

	viper.SetConfigName("ConfData") // name of config file (without extension)
	viper.SetConfigType("json")
	viper.AddConfigPath("./")


	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	//是否开启文件监控
	//viper.WatchConfig()
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	fmt.Println("Config file changed:", in.Name)
	//
	//	UpdateConfAll()
	//})

	UpdateConfAll()
}