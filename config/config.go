package config

import (
	"errors"
	"fmt"
	"github.com/Unknwon/goconfig"
	"os"
	"path/filepath"
)

const configFile = "/conf/conf.ini"

var File *goconfig.ConfigFile

// 加载此文件时候会走初始化
func init() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	println(currentDir)
	configPath := filepath.Join(currentDir, configFile)
	// configPath := currentDir + configFile
	println(configPath)

	// 配置文件不存在，直接就停在这里不往下走了。
	if !fileExists(configPath) {
		panic(errors.New("config file not exists"))
	}

	// 增加健壮性，在控制执行读取配置时，可以自定配置路径
	len := len(os.Args)
	if len > 1 {
		dir := os.Args[1]
		if dir != "" {
			configPath = filepath.Join(dir, configFile)
		}
	}

	File, err = goconfig.LoadConfigFile(configPath)
	if err != nil {
		// panic(err)
		fmt.Println("读取配置文件失败", err.Error())
		return
	}

	// fmt.Println(File)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)

}

func AAA() {
	fmt.Println("aaaaaa")
}
