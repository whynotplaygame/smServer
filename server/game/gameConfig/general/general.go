package general

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"path/filepath"
)

type general struct {
	Title            string          `json:"title"`
	GArr             []generalDetail `json:"list"`
	GMap             map[int]generalDetail
	totalProbability int
}
type generalDetail struct {
	Name         string `json:"name"`
	CfgId        int    `json:"cfgId"`
	Force        int    `json:"force"`    //武力
	Strategy     int    `json:"strategy"` //策略
	Defense      int    `json:"defense"`  //防御
	Speed        int    `json:"speed"`    //速度
	Destroy      int    `json:"destroy"`  //破坏力
	ForceGrow    int    `json:"force_grow"`
	StrategyGrow int    `json:"strategy_grow"`
	DefenseGrow  int    `json:"defense_grow"`
	SpeedGrow    int    `json:"speed_grow"`
	DestroyGrow  int    `json:"destroy_grow"`
	Cost         int8   `json:"cost"`
	Probability  int    `json:"probability"`
	Star         int8   `json:"star"`
	Arms         []int  `json:"arms"`
	Camp         int8   `json:"camp"`
}

var General = &general{}
var generalFile = "/conf/game/general/general.json"

func (g *general) Load() {
	g.GMap = make(map[int]generalDetail) // 切片使用前要初始化
	g.GArr = make([]generalDetail, 0)

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	println(currentDir)
	configFile := filepath.Join(currentDir, generalFile)
	// configPath := filepath.Join(currentDir, facilityPath)
	// configPath := currentDir + basicFile
	println(configFile)

	// 增加健壮性，在控制执行读取配置时，可以自定配置路径
	length := len(os.Args)
	if length > 1 {
		dir := os.Args[1]
		if dir != "" {
			configFile = filepath.Join(currentDir, generalFile)
		}
	}

	data, err := os.ReadFile(configFile) // 教程里的ioutil.read已被弃用。官网推荐使用os.readfile
	if err != nil {
		log.Println("读取文件失败：", err)
		panic(err)
	}
	err = json.Unmarshal(data, g)
	if err != nil {
		log.Println("解析json错误：", err)
	}

	for _, v := range g.GArr {
		g.GMap[v.CfgId] = v
		g.totalProbability += v.Probability
	}

	log.Println("加载城池设施配置成功")
}

// 随机武将
func (g *general) Rand() int {
	// 7+12 = 19,  0-19    8
	rate := rand.Intn(g.totalProbability)
	var cur = 0
	for _, v := range g.GArr {
		if rate >= cur && rate < cur+v.Probability {
			return v.CfgId
		}
		cur += v.Probability
	}
	return 0
}
