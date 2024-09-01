package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() (err error){
	//使用配置文件的方式
	viper.SetConfigFile("./settings/config.yaml")

	if err = viper.ReadInConfig();err != nil{
		fmt.Printf("viper.ReadInConfig() error : %v",err)
		return
	}

	//监控配置文件是否发生变化
	viper.WatchConfig()
	//在改变后使用回调函数，在下次访问时使用改变后的值
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件被修改，请重启系统")
	})

	return
}