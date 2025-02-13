package gameConfig

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

const (
	TypeDurable        = 1 //耐久
	TypeCost           = 2
	TypeArmyTeams      = 3 //队伍数量
	TypeSpeed          = 4 //速度
	TypeDefense        = 5 //防御
	TypeStrategy       = 6 //谋略
	TypeForce          = 7 //攻击武力
	TypeConscriptTime  = 8 //征兵时间
	TypeReserveLimit   = 9 //预备役上限
	TypeUnkonw         = 10
	TypeHanAddition    = 11
	TypeQunAddition    = 12
	TypeWeiAddition    = 13
	TypeShuAddition    = 14
	TypeWuAddition     = 15
	TypeDealTaxRate    = 16 //交易税率
	TypeWood           = 17
	TypeIron           = 18
	TypeGrain          = 19
	TypeStone          = 20
	TypeTax            = 21 //税收
	TypeExtendTimes    = 22 //扩建次数
	TypeWarehouseLimit = 23 //仓库容量
	TypeSoldierLimit   = 24 //带兵数量
	TypeVanguardLimit  = 25 //前锋数量
)

type conditions struct {
	Type  int `json:"type"`
	Level int `json:"level"`
}

type facility struct {
	Title      string       `json:"title"`
	Des        string       `json:"des"`
	Name       string       `json:"name"`
	Type       int8         `json:"type"`
	Additions  []int8       `json:"additions"`
	Conditions []conditions `json:"conditions"`
	Levels     []fLevel     `json:"levels"`
}

type NeedRes struct {
	Decree int `json:"decree"`
	Grain  int `json:"grain"`
	Wood   int `json:"wood"`
	Iron   int `json:"iron"`
	Stone  int `json:"stone"`
	Gold   int `json:"gold"`
}

type fLevel struct {
	Level  int     `json:"level"`
	Values []int   `json:"values"`
	Need   NeedRes `json:"need"`
	Time   int     `json:"time"` //升级需要的时间
}

type conf struct {
	Name string
	Type int8
}

type facilityConf struct {
	Title     string `json:"title"`
	List      []conf `json:"list"`
	facilitys map[int8]*facility
}

var FaiclityConfig = facilityConf{}

const facilityFile = "/conf/game/facility/facility.json"
const facilityPath = "/conf/game/facility/"

func (f *facilityConf) Load() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	println(currentDir)
	configFile := filepath.Join(currentDir, facilityFile)
	configPath := filepath.Join(currentDir, facilityPath)
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
			configFile = filepath.Join(currentDir, facilityFile)
			configPath = filepath.Join(currentDir, facilityPath)
		}
	}

	data, err := os.ReadFile(configFile) // 教程里的ioutil.read已被弃用。官网推荐使用os.readfile
	if err != nil {
		log.Println("读取文件失败：", err)
		panic(err)
	}
	err = json.Unmarshal(data, f)
	if err != nil {
		log.Println("解析json错误：", err)
	}

	f.facilitys = make(map[int8]*facility, len(f.List))

	fils, err := os.ReadDir(configPath)
	if err != nil {
		log.Println("读取设施文件夹失败：", err)
		panic(err)
	}
	for _, file := range fils {
		if file.IsDir() {
			continue
		}
		if file.Name() == "facility" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(configPath, file.Name()))
		if err != nil {
			log.Println("读取设施文件失败：", err)
			panic(err)
		}

		fac := &facility{}
		err = json.Unmarshal(data, fac)
		if err != nil {
			log.Println("解析设施json数据失败：", err)
			panic(err)
		}
		f.facilitys[fac.Type] = fac
	}
	log.Println("加载城池设施配置成功")
}
