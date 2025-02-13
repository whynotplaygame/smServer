package gameConfig

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"smServer/server/game/global"
)

type mapData struct {
	Width  int     `json:"w"`
	Height int     `json:"h"`
	List   [][]int `json:"list"`
}

// 加载地图上每个单元格的属性,定义文件放在了资料的map.json中
// 要把 mapData的list中每个值，转换成nationmap中的值
type NationalMap struct {
	MId   int  `xorm:"mid"`
	X     int  `xorm:"x"`
	Y     int  `xorm:"y"`
	Type  int8 `xorm:"type"`
	Level int8 `xorm:"level"`
}

const (
	MapBuildSysFortress = 50 //系统要塞
	MapBuildSysCity     = 51 //系统城市
	MapBuildFortress    = 56 //玩家要塞
)

type mapRes struct {
	Confs    map[int]NationalMap
	SysBuild map[int]NationalMap // 存放系统要塞和系统城市
}

var MapRes = &mapRes{
	Confs:    make(map[int]NationalMap),
	SysBuild: make(map[int]NationalMap),
}

const mapFile = "/conf/game/map.json"

func (m *mapRes) Load() {

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	println(currentDir)
	configPath := filepath.Join(currentDir, mapFile)
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
	mapData := &mapData{}
	err = json.Unmarshal(data, mapData)
	if err != nil {
		log.Println("解析json错误：", err)
	}
	global.MapWith = mapData.Width
	global.MapHeight = mapData.Height
	log.Println("map list len:", len(mapData.List))

	for index, v := range mapData.List {
		t := int8(v[0])
		l := int8(v[1])
		nm := NationalMap{
			X:     index % global.MapWith,
			Y:     index / global.MapWith,
			Type:  t,
			Level: l,
			MId:   index,
		}
		m.Confs[index] = nm
		if t == MapBuildSysCity || t == MapBuildFortress {
			m.SysBuild[index] = nm
		}
	}
	log.Println("加载地图配置成功")
}
