package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"net/http"
	"os"
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

	iqiyiScore     string = "http://pcw-api.iqiyi.com/video/score/getsnsscore?qipu_ids=230419701&tvid=230419701&pageNo=1"
	iqiyiPlayTimes string = "https://pcw-api.iqiyi.com/video/video/hotplaytimes/230419701"
	iqiyiRank string = "http://top.iqiyi.com/shaoer.html"

	tencentScore string = "https://v.qq.com/x/cover/to61xna5r970zmo/e0027wpnpye.html"
	tencentRank  string = "https://v.qq.com/x/hotlist/search/?channel=106"

	pptvScore             string = "http://v.pptv.com/page/JWdQzjacDEqtK5M.html?spm=v_show_web.0.1.3.1.3.1.3.2.1"
	pptvRank              string = "http://top.pptv.com/kid?fb=1"
	pptvPlayTimesBySeries string = "http://v.pptv.com/show/JWdQzjacDEqtK5M.html?spm=pc_top_web.0.1.2.0.2.2.0.7.1.0"

	sheetNameOne   string = "新媒体播放数据对比增幅"
	sheetNameTwo   string = "芒果分集播放量"
	sheetNameThree string = "PP分集播放量"
	sheetNameFour  string = "优酷分集播放量"
)

var (
	h           bool
	p           string
	currentDate = time.Now().Format("2006-01-02 00:00:00")
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

var fillSeq = [5]string{"iqiyi", "tencent", "mgtv", "pptv", "youku"}
var dataSeq = make(map[string]Collector)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&p, "p", "collect.xlsx", "设置excel文件的路径")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `
Usage: fetch [excel]

Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	mgtv := getMGTVData()
	iqiyi := getIQiyiData()
	tencent := getTencentData()
	pptv := getPPTVData()
	dataSeq["mgtv"] = *mgtv
	dataSeq["iqiyi"] = *iqiyi
	dataSeq["tencent"] = *tencent
	dataSeq["pptv"] = *pptv

	fillExcel()
}

func fillSheetOneData(platform string, row *xlsx.Row, startIndex int) {
	if startIndex < 0 {
		panic("Please check current date column exists in sheet one")
	}
	cells := row.Cells
	fmt.Printf("%s data as below: \n times: %s \n score: %s \n rank: %s\n serial: %s\n", platform, dataSeq[platform].playTimes, dataSeq[platform].score, dataSeq[platform].rank, dataSeq[platform].seriesTimes)
	style := cells[startIndex].GetStyle()
	style.Font.Size = 9
	cells[startIndex].Value = dataSeq[platform].playTimes
	cells[startIndex+1].Value = dataSeq[platform].score
	cells[startIndex+2].Value = dataSeq[platform].rank
}

func fillSheetsData(platform string, sheet *xlsx.Sheet) {
	startIndex := -1
	for rowIndex, row := range sheet.Rows {
		if rowIndex > 53 {
			break
		}
		if rowIndex == 1 {
			startIndex = getStartIndexByMatchDate(row)
			if startIndex < 0 {
				fmt.Println("Please check current date column exists in sheet two: " + platform)
			}
		} else if rowIndex > 1 && startIndex > 0 {
			//第一列的值正好是第几集，1，2，3
			value, err := strconv.ParseFloat(dataSeq[platform].seriesTimes[row.Cells[0].Value], 64)
			if err != nil {
				fmt.Println("Format float err when parse " + platform + " on sheet two")
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

func fillExcel() {
	xlFile, err := xlsx.OpenFile(p)
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
			iqiyiStartIndex := -1
			for rowIndex, row := range sheet.Rows {
				if rowIndex == 1 {
					//获取填充的列位置
					startIndex = getStartIndexByMatchDate(row)
				}
				if rowIndex == 3 {
					//填充爱奇艺的数据
					fillSheetOneData("iqiyi", row, startIndex)
				}
				if rowIndex == 4 {
					//填充腾讯的数据
					fillSheetOneData("tencent", row, startIndex)
				}
				if rowIndex == 5 {
					//填充芒果TV的数据
					fillSheetOneData("mgtv", row, startIndex)
				}
				if rowIndex == 6 {
					//填充PPTV的数据
					fillSheetOneData("pptv", row, startIndex)
				}
				if rowIndex == 11 {
					//获取爱奇艺趋势图位置
					for j, cell := range row.Cells {
						if j == 0 {
							continue
						}
						switch cell.Type() {
						case xlsx.CellTypeNumeric:
							t, _ := cell.GetTime(false)
							tstr := t.Format("2006-01-02 00:00:00")
							if strings.EqualFold(tstr, currentDate) {
								iqiyiStartIndex = j
							}
						}
					}
					if iqiyiStartIndex < 0 {
						panic("爱奇异趋势图的当前日期不存在")
					}
				}
				if rowIndex == 12 {
					cells := row.Cells
					style := cells[iqiyiStartIndex].GetStyle()
					style.Font.Size = 9
					cells[iqiyiStartIndex].Value = dataSeq["iqiyi"].playTimes

				}
			}
		}
		if strings.EqualFold(sheet.Name, sheetNameTwo) {
			fillSheetsData("mgtv", sheet)
		}
		if strings.EqualFold(sheet.Name, sheetNameThree) {
			fillSheetsData("pptv", sheet)
		}
		if strings.EqualFold(sheet.Name, sheetNameFour) {
			fillSheetsData("youku", sheet)
		}
	}
	xlFile.Save(p)
}

func getStartIndexByMatchDate(row *xlsx.Row) int {
	var startIndex = -1
	for j, cell := range row.Cells {
		if j == 0 {
			//fmt.Printf("\n")
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

func getPPTVData() *Collector {
	c := new(Collector)

	c.platform = "pptv"
	//获取排名
	doc, err := goquery.NewDocument(pptvRank)
	doc.Find("body").Find("ul.cf").Find("li").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a[title=\"宇宙护卫队\"]")
		if s != nil && strings.TrimSpace(s.Text()) == "宇宙护卫队" {
			span := selection.Find("span")
			rank, err := strconv.ParseInt(span.Text(), 10, 64)
			if err != nil {
				fmt.Println("[Error]Get PPTV Rank fail")
			}
			c.rank = strconv.Itoa(int(rank))
		}
	})

	//获取评分
	resp, _ := http.Get(pptvScore)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	pat := "<b class=\"score\">[0-9]+\\.?[0-9]*</b>"
	reg, _ := regexp.Compile(pat)
	span := reg.Find(body)
	c.score = GetBetweenStr(string(span), ">", "<")
	fmt.Println("Get PPTV Score: " + c.score)

	//获取播放
	pat = "<li>播放：[0-9]+\\.?[0-9]*万"
	reg, _ = regexp.Compile(pat)
	span = reg.Find(body)
	tmp := GetBetweenStr(string(span), "<li>播放：", "万")
	c.playTimes = strings.Split(tmp, "：")[1]

	resp, _ = http.Get(pptvPlayTimesBySeries)
	body, _ = ioutil.ReadAll(resp.Body)
	jsonStr := GetBetweenStr(string(body), "var webcfg =", "\n")
	jsonStr = GetBetweenStr(jsonStr, "=", ";")

	//fmt.Println(jsonStr)
	res, err := simplejson.NewJson([]byte(jsonStr))
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	rows, err := res.Get("playList").Get("data").Get("list").Array()
	c.seriesTimes = make(map[string]string)
	for _, row := range rows {
		if each_map, ok := row.(map[string]interface{}); ok {
			if rank, ok := each_map["rank"].(json.Number); ok {
				rank_int, err := strconv.ParseInt(string(rank), 10, 0)
				rank_int = rank_int + 1
				if err != nil {
					panic(err)
				}
				c.seriesTimes[strconv.FormatInt(rank_int, 10)] = strings.Split(each_map["pv"].(string), "万")[0]
			}
		}
	}

	return c
}

func getTencentData() *Collector {
	c := new(Collector)

	c.platform = "tencent"
	//获取评分
	resp, _ := http.Get(tencentScore)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	pat := "\"score\":\"[0-9]+\\.?[0-9]*"
	reg, _ := regexp.Compile(pat)
	span := reg.Find(body)
	c.score = strings.Split(string(span), "\":\"")[1]
	fmt.Println("Get Tencent Score: " + c.score)

	//获取排名
	doc, err := goquery.NewDocument(tencentRank)
	doc.Find("body").Find("ul.table_list._cont").Find("li").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a[title=\"宇宙护卫队\"]")
		if s != nil && strings.TrimSpace(s.Text()) == "宇宙护卫队" {
			span := selection.Find("span")
			c.rank = span.Text()
		}
	})
	//获取播放次数
	pat = "<em id=\"mod_cover_playnum\" class=\"num\">[0-9]+\\.?[0-9]*亿</em>"
	reg, _ = regexp.Compile(pat)
	span = reg.Find(body)
	tmp := GetBetweenStr(string(span), ">", "亿</em>")
	tmpFloat, err := strconv.ParseFloat(tmp, 64)
	if err != nil {
		fmt.Println("Tencent Playtime format err: ")
		panic(err)
	}
	tmpInt := tmpFloat * 10000
	c.playTimes = strconv.Itoa(int(tmpInt))
	fmt.Println("Get Tencent PlayTimes: " + c.playTimes)

	return c
}

func getMGTVData() *Collector {
	c := new(Collector)

	c.platform = "mgtv"

	//获取评分
	resp, _ := http.Get(mgtvScore)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	pat := "<span class=\"score\">[0-9]+\\.?[0-9]*</span>"
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

type IQiyiScore struct {
	Data []IQiyiScoreData `json:"data"`
}

type IQiyiScoreData struct {
	ID    int64   `json:"qipu_id"`
	Score float64 `json:"sns_score"`
}

type IQiyiPlayTimes struct {
	Data []IQiyiPlayTimesData `json:"data"`
}

type IQiyiPlayTimesData struct {
	Hot int `json:"hot"`
}

func getIQiyiData() *Collector {
	c := new(Collector)
	c.platform = "iqiyi"

	//获取排名
	doc, err := goquery.NewDocument(iqiyiRank)
	doc.Find("body").Find("ul.topDetails-list").Find("li").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a[title=\"宇宙护卫队\"]")
		if s != nil && strings.TrimSpace(s.Text()) == "宇宙护卫队" {
			span := selection.Find("i.array")
			if err != nil {
				fmt.Println("[Error]Get IQIYI Rank fail")
			}
			c.rank = span.Text()
		}
	})
	//获取评分
	resp, _ := http.Get(iqiyiScore)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var iqiyiScore IQiyiScore
	if err := json.Unmarshal(body, &iqiyiScore); err != nil {
		fmt.Println("================IQiyiScore json str 转struct==")
		fmt.Println(err)
	}
	c.score = strconv.FormatFloat(iqiyiScore.Data[0].Score, 'f', 1, 64)

	//获取热度
	resp, _ = http.Get(iqiyiPlayTimes)
	body, _ = ioutil.ReadAll(resp.Body)
	var iqiyiPlayTimes IQiyiPlayTimes
	if err := json.Unmarshal(body, &iqiyiPlayTimes); err != nil {
		fmt.Println("================IQiyiPlayTimes json str 转struct==")
		fmt.Println(err)
	}
	c.playTimes = strconv.Itoa(iqiyiPlayTimes.Data[0].Hot)

	return c
}
