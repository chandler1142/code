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
	platform  string
	score     string
	rank      string
	playTimes string
}

const (
	cartoonName   string = "宇宙护卫队"
	mgtvScore     string = "https://so.mgtv.com/so/k-%E5%AE%87%E5%AE%99%E6%8A%A4%E5%8D%AB%E9%98%9F"
	mgtvRank      string = "https://rc.mgtv.com/pc/ranklist?&c=50&t=day&limit=30&rt=c&t=%s"
	mgtvPlayTimes string = "https://vc.mgtv.com/v2/dynamicinfo?_support=10000000&cid=326647&_=%s"

	sheetNameOne   string = "新媒体播放数据对比增幅"
	sheetNameTwo   string = "芒果分集播放量"
	sheetNameThree string = "PP分集播放量"

	path string = "pingfen/collect.xlsx"
)

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
	fmt.Println(mgtv.score)
	dataSeq["mgtv"] = *mgtv

	fillExcel()
}

func fillExcelSheetOne() {

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
	//for _, sheet := range xlFile.Sheets {
	//	//fmt.Printf("Sheet Name: %s\n", sheet.Name)
	//	if strings.EqualFold(sheet.Name, sheetNameOne) {
	//		startIndex := -1
	//		currentDate := time.Now().Format("2006-01-02")
	//		fmt.Println("currentData: " + currentDate)
	//		for rowIndex, row := range sheet.Rows {
	//			if rowIndex == 1 {
	//				for j, cell := range row.Cells {
	//					if j == 0 {
	//						fmt.Printf("\n")
	//					}
	//					switch cell.Type() {
	//					case xlsx.CellTypeString:
	//						fmt.Printf("str %d %s\t", j, cell.String())
	//					case xlsx.CellTypeStringFormula:
	//						fmt.Printf("formula %d %s\t", j, cell.Formula())
	//					case xlsx.CellTypeNumeric:
	//						//x, _ := cell.Int64()
	//						//fmt.Printf("int %d %d\t", j, x)
	//						t, _ := cell.GetTime(false)
	//						//fmt.Printf("date %d %v\t", j, t)
	//						tstr := t.Format("2006-01-02")
	//						if strings.EqualFold(tstr, currentDate) {
	//							startIndex = j
	//						}
	//					case xlsx.CellTypeBool:
	//						fmt.Printf("bool %d %v\t", j, cell.Bool())
	//					case xlsx.CellTypeDate:
	//						//t, _ := cell.GetTime(false)
	//						//fmt.Printf("date %d %v\t", j, t)
	//						//break
	//						t, _ := cell.GetTime(false)
	//						tstr := t.Format("2006-01-02")
	//						if strings.EqualFold(tstr, currentDate) {
	//							startIndex = j
	//						}
	//					}
	//				}
	//			}
	//			//}
	//			if rowIndex == 3 {
	//				//填充爱奇艺的数据
	//			}
	//			if rowIndex == 5 {
	//				//填充芒果TV的数据
	//				cells := row.Cells
	//				totalLength := len(cells)
	//				fmt.Printf("totalLength: %d \n", totalLength)
	//				fmt.Printf("startIndex: %d \n", startIndex)
	//				fmt.Printf("mgtv data as below: \n times: %s \n score: %s \n rank: %s\n", dataSeq["mgtv"].playTimes, dataSeq["mgtv"].score, dataSeq["mgtv"].rank)
	//				cells[startIndex].Value = dataSeq["mgtv"].playTimes
	//				cells[startIndex+1].Value = dataSeq["mgtv"].score
	//				cells[startIndex+2].Value = dataSeq["mgtv"].rank
	//			}
	//		}
	//
	//	}
	//
	//}
	xlFile.Save(path)
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
	fmt.Printf("Get MGTV rank data, length: %d", len(mgtvRankData.Data))
	for _, data := range mgtvRankData.Data {
		if strings.Contains(data.Name, cartoonName) {
			c.rank = string(data.VideoIndex)
			break
		}
	}
	if c.rank == "" {
		c.rank = ""
	}
	fmt.Println("Get MGTV Rank: " + c.rank)

	//获取播放次数
	resp, _ = http.Get(fmt.Sprintf(mgtvPlayTimes, strconv.FormatInt(second, 10)))
	body, _ = ioutil.ReadAll(resp.Body)
	var mgtvPlayTimes MgtvPlayTimes
	if err := json.Unmarshal(body, &mgtvPlayTimes); err != nil {
		fmt.Println("================mgtvPlayTimes json str 转struct==")
		fmt.Println(err)
	}

	playTimesFloat,err := strconv.ParseFloat(strconv.Itoa(mgtvPlayTimes.Data.All/10000) + "." + strconv.Itoa(mgtvPlayTimes.Data.All%10000),10)
	if err != nil {
		fmt.Println("parse float err")
		panic(err)
	}
	c.playTimes = strconv.FormatFloat(playTimesFloat, 'f', 2, 64)

	return c
}
