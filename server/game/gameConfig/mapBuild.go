package gameConfig

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

type cfg struct {
	Type     int8   `json:"type"`
	Name     string `json:"name"`
	Level    int8   `json:"level"`
	Grain    int    `json:"grain"`
	Wood     int    `json:"wood"`
	Iron     int    `json:"iron"`
	Stone    int    `json:"stone"`
	Durable  int    `json:"durable"`
	Defender int    `json:"defender"`
}

type mapBuildConf struct {
	Title  string `json:"title"`
	Cfg    []cfg  `json:"cfg"`
	cfgMap map[int8][]cfg
}

var MapBuildConf = &mapBuildConf{
	cfgMap: make(map[int8][]cfg),
}

const mapBuilderConfFile = "/conf/game/map_Build.json"

func (m *mapBuildConf) Load() {

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	println(currentDir)
	configPath := filepath.Join(currentDir, mapBuilderConfFile)
	// configPath := currentDir + basicFile
	println(configPath)

	// 配置文件不存在，直接就停在这里不往下走了。
	if !fileExists(configPath) {
		panic(errors.New("config file not exists"))
	}

	// 增加健壮性，在控制执行读取配置时，可以自定配置路径
	length := len(os.Args)
	if length > 1 {
		dir := os.Args[1]
		if dir != "" {
			configPath = filepath.Join(dir, mapBuilderConfFile)
		}
	}

	data, err := os.ReadFile(configPath) // 教程里的ioutil.read已被弃用。官网推荐使用os.readfile
	if err != nil {
		log.Println("读取文件失败：", err)
		panic(err)
	}
	err = json.Unmarshal(data, m)
	if err != nil {
		log.Println("解析json错误：", err)
	}

	for _, v := range m.Cfg {
		_, ok := m.cfgMap[v.Type]
		if !ok {
			m.cfgMap[v.Type] = make([]cfg, 0)
		} else {
			m.cfgMap[v.Type] = append(m.cfgMap[v.Type], v)
		}
	}

	log.Println("加载城池建筑配置成功")
}

func (m *mapBuildConf) Buildconfig(buildType int8, level int8) *cfg {
	cfgs := m.cfgMap[buildType]
	for _, v := range cfgs {
		if v.Level == level {
			return &v
		}
	}
	return nil
}
