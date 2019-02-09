package main

import (
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Collector struct {
	platform    string
	score       string
	rank        string
	playTimes   string
	seriesTimes map[string]string
}

const (
	cartoonName string = "宇宙护卫队"

	mgtvScore             string = "https://so.mgtv.com/so/k-%E5%AE%87%E5%AE%99%E6%8A%A4%E5%8D%AB%E9%98%9F"
	mgtvRank              string = "https://rc.mgtv.com/pc/ranklist?&c=50&t=day&limit=30&rt=c&t=%s"
	mgtvPlayTimes         string = "https://vc.mgtv.com/v2/dynamicinfo?_support=10000000&cid=326647&_=%s"
	mgtvPlayTimesBySeries string = "https://pcweb.api.mgtv.com/episode/list?video_id=4619079&page=%d&size=25&cxid=&version=5.5.35&_support=10000000&_=%s"

	iqiyiScore string = "http://www.iqiyi.com/a_19rrh51i8p.html?vfm=2008_aldbd"

	sheetNameOne   string = "新媒体播放数据对比增幅"
	sheetNameTwo   string = "芒果分集播放量"
	sheetNameThree string = "PP分集播放量"

	path string = "pingfen/collect.xlsx"
)

var currentDate = time.Now().Format("2006-01-02 00:00:00")

func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n+1:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

var fillSeq = [5]string{"iqiyi", "tecent", "mgtv", "pptv", "youku"}
var dataSeq = make(map[string]Collector)

func main() {
	mgtv := getMGTVData()
	dataSeq["mgtv"] = *mgtv
	fillExcel()
}

func fillExcel() {
	excelName := path
	xlFile, err := xlsx.OpenFile(excelName)
	if xlFile == nil {
		fmt.Println("No such excel file named collect.xlsx")
		return
	}
	if err != nil {
		fmt.Printf("open failed: %s\n", err)
	}

	for _, sheet := range xlFile.Sheets {
		fmt.Printf("\n=================== Start to process %s ========================== \n", sheet.Name)
		if strings.EqualFold(sheet.Name, sheetNameOne) {
			startIndex := -1
			for rowIndex, row := range sheet.Rows {
				if rowIndex == 1 {
					//获取填充的列位置
					startIndex = getStartIndexByMatchDate(row)
				}
				if rowIndex == 3 {
					//填充爱奇艺的数据
				}
				if rowIndex == 5 {
					//填充芒果TV的数据
					cells := row.Cells
					//fmt.Printf("startIndex: %d \n", startIndex)
					fmt.Printf("MGTV data as below: \n times: %s \n score: %s \n rank: %s\n serial: %s\n", dataSeq["mgtv"].playTimes, dataSeq["mgtv"].score, dataSeq["mgtv"].rank, dataSeq["mgtv"].seriesTimes)
					cells[startIndex].Value = dataSeq["mgtv"].playTimes
					cells[startIndex+1].Value = dataSeq["mgtv"].score
					cells[startIndex+2].Value = dataSeq["mgtv"].rank
				}
			}

		}

		if strings.EqualFold(sheet.Name, sheetNameTwo) {
			startIndex := -1
			for rowIndex, row := range sheet.Rows {
				if rowIndex > 53 {
					break
				}
				if rowIndex == 1 {
					startIndex = getStartIndexByMatchDate(row)
				} else if rowIndex > 1 {
					//第一列的值正好是第几集，1，2，3
					value, err := strconv.ParseFloat(dataSeq["mgtv"].seriesTimes[row.Cells[0].Value], 64)
					if err != nil {
						fmt.Println("Format float err when parse MGTV on sheet two")
					}
					if startIndex > len(row.Cells)-1 {
						for i := 0; startIndex >= len(row.Cells)-1; i++ {
							row.AddCell()
						}
					}
					row.Cells[startIndex].SetFloat(value)
				}
			}
		}
	}
	xlFile.Save(path)
}

func getStartIndexByMatchDate(row *xlsx.Row) int {
	var startIndex = -1
	for j, cell := range row.Cells {
		if j == 0 {
			fmt.Printf("\n")
		}
		switch cell.Type() {
		case xlsx.CellTypeNumeric:
			t, _ := cell.GetTime(false)
			tstr := t.Format("2006-01-02 00:00:00")
			if strings.EqualFold(tstr, currentDate) {
				startIndex = j
			}
			//后面几个sheet强制设置一下类型
			cell.NumFmt = "m\"月\"d\"日\";@"
		}
	}
	return startIndex
}

type MgtvRankData struct {
	Data []MgtvRank `json:"data"`
}

type MgtvRank struct {
	VideoIndex int    `json:"videoIndex"`
	Name       string `json:"name"`
}

type MgtvPlayTimes struct {
	Data MgtvPlayData `json:"data"`
}

type MgtvPlayData struct {
	All int `json:"all"`
}

type MgtvPlaySerial struct {
	Code int                `json:"code"`
	Data MgtvPlaySerialData `json: "data"`
}

type MgtvPlaySerialData struct {
	DataList []MgtvPlaySerialDataElement `json:"list"`
}

type MgtvPlaySerialDataElement struct {
	T1    string `json:"t1"`
	Count string `json:"playcnt"`
}

func getMGTVData() *Collector {
	c := new(Collector)

	c.platform = "mgtv"

	//获取评分
	resp, _ := http.Get(mgtvScore)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	pat := "<span class=\"score\">[0-9]+\\.[0-9]+</span>"
	reg, _ := regexp.Compile(pat)
	span := reg.Find(body)
	c.score = GetBetweenStr(string(span), ">", "<")
	fmt.Println("Get MGTV Score: " + c.score)

	//获取排名
	second := time.Now().Unix()
	resp, _ = http.Get(fmt.Sprintf(mgtvRank, strconv.FormatInt(second, 10)))
	body, _ = ioutil.ReadAll(resp.Body)
	var mgtvRankData MgtvRankData
	if err := json.Unmarshal(body, &mgtvRankData); err != nil {
		fmt.Println("================mgtvRankData json str 转struct==")
		fmt.Println(err)
	}
	fmt.Printf("Get MGTV rank data, length: %d \n", len(mgtvRankData.Data))
	for _, data := range mgtvRankData.Data {
		if strings.Contains(data.Name, cartoonName) {
			c.rank = string(data.VideoIndex)
			break
		}
	}
	if c.rank == "" {
		c.rank = ""
	}

	//获取播放次数
	resp, _ = http.Get(fmt.Sprintf(mgtvPlayTimes, strconv.FormatInt(second, 10)))
	body, _ = ioutil.ReadAll(resp.Body)
	var mgtvPlayTimes MgtvPlayTimes
	if err := json.Unmarshal(body, &mgtvPlayTimes); err != nil {
		fmt.Println("================mgtvPlayTimes json str 转struct==")
		fmt.Println(err)
	}

	playTimesFloat, err := strconv.ParseFloat(strconv.Itoa(mgtvPlayTimes.Data.All/10000)+"."+strconv.Itoa(mgtvPlayTimes.Data.All%10000), 10)
	if err != nil {
		fmt.Println("parse float err")
		panic(err)
	}
	c.playTimes = strconv.FormatFloat(playTimesFloat, 'f', 2, 64)

	c.seriesTimes = make(map[string]string)
	//一共三页
	for i := 1; i <= 3; i++ {
		resp, _ = http.Get(fmt.Sprintf(mgtvPlayTimesBySeries, i, second))
		body, _ = ioutil.ReadAll(resp.Body)
		var mgtvPlaySerial MgtvPlaySerial
		if err := json.Unmarshal(body, &mgtvPlaySerial); err != nil {
			fmt.Println("================mgtvPlaySerial json str 转struct==")
			fmt.Println(err)
		}

		for _, element := range mgtvPlaySerial.Data.DataList {
			c.seriesTimes[element.T1] = strings.Split(element.Count, "万")[0]
		}
	}

	return c
}
