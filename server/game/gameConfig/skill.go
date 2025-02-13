package gameConfig

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var Skill skill

type skill struct {
	skills       []Conf
	skillConfMap map[int]Conf
	outline      outline
}

type trigger struct {
	Type int    `json:"type"`
	Des  string `json:"des"`
}

type triggerType struct {
	Des  string    `json:"des"`
	List []trigger `json:"list"`
}

type effect struct {
	Type   int    `json:"type"`
	Des    string `json:"des"`
	IsRate bool   `json:"isRate"`
}

type effectType struct {
	Des  string   `json:"des"`
	List []effect `json:"list"`
}

type target struct {
	Type int    `json:"type"`
	Des  string `json:"des"`
}

type targetType struct {
	Des  string   `json:"des"`
	List []target `json:"list"`
}

type outline struct {
	TriggerType triggerType `json:"trigger_type"` //触发类型
	EffectType  effectType  `json:"effect_type"`  //效果类型
	TargetType  targetType  `json:"target_type"`  //目标类型
}

type level struct {
	Probability int   `json:"probability"`  //发动概率
	EffectValue []int `json:"effect_value"` //效果值
	EffectRound []int `json:"effect_round"` //效果持续回合数
}

type Conf struct {
	CfgId         int     `json:"cfgId"`
	Name          string  `json:"name"`
	Trigger       int     `json:"trigger"` //发起类型
	Target        int     `json:"target"`  //目标类型
	Des           string  `json:"des"`
	Limit         int     `json:"limit"`          //可以被武将装备上限
	Arms          []int   `json:"arms"`           //可以装备的兵种
	IncludeEffect []int   `json:"include_effect"` //技能包括的效果
	Levels        []level `json:"levels"`
}

const skillFile = "/conf/game/skill/skill_outline.json"
const skillPath = "/conf/game/skill/"

func (s *skill) Load() {
	s.skills = make([]Conf, 0)
	s.skillConfMap = make(map[int]Conf)

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	println(currentDir)
	configFile := filepath.Join(currentDir, skillFile)
	configPath := filepath.Join(currentDir, skillPath)
	// configPath := currentDir + basicFile
	println(configPath)

	// 增加健壮性，在控制执行读取配置时，可以自定配置路径
	length := len(os.Args)
	if length > 1 {
		dir := os.Args[1]
		if dir != "" {
			configFile = filepath.Join(currentDir, skillFile)
			configPath = filepath.Join(currentDir, skillPath)
		}
	}

	data, err := os.ReadFile(configFile) // 教程里的ioutil.read已被弃用。官网推荐使用os.readfile
	if err != nil {
		log.Println("读取文件失败：", err)
		panic(err)
	}
	err = json.Unmarshal(data, &s.outline)
	if err != nil {
		log.Println("技能json读取文件失败：", err)
		panic(err)
	}

	fils, err := os.ReadDir(configPath)
	if err != nil {
		log.Println("读取技能文件夹失败：", err)
		panic(err)
	}
	for _, v := range fils {
		if v.IsDir() {
			name := v.Name()
			dirFile := filepath.Join(configPath, name)
			skillFiles, err := os.ReadDir(dirFile)
			if err != nil {
				log.Println(name + "技能文件读取失败")
				panic(err)
			}
			for _, sv := range skillFiles {
				if sv.IsDir() {
					continue
				}
				fileJson := filepath.Join(dirFile, sv.Name())
				conf := Conf{}
				data, err := os.ReadFile(fileJson)
				if err != nil {
					log.Println(name + "技能文件格式错误")
					panic(err)
				}
				err = json.Unmarshal(data, &conf)
				if err != nil {
					log.Println(name + "技能文件格式错误")
					panic(err)
				}
				s.skills = append(s.skills, conf)
				s.skillConfMap[conf.CfgId] = conf
			}
		}

	}
	log.Println("s:", s)
	log.Println("加载技能配置成功")
}
