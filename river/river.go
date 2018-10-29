package river

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// WaterLevel 水位記録の１行
type WaterLevel struct {
	Time  time.Time
	Level float64
}

// River 河川情報
type River struct {
	Name        string //河川名
	WaterLevels []WaterLevel
}

// Download 河川の詳細ページから表の内容をダウンロードする
func Download(path string) {
	fmt.Println(path)
	parm, _ := url.Parse(path)

	res, err := http.Get(buildURL(path))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	river := River{
		Name:        getRiverName(doc),
		WaterLevels: getWaterLevels(doc),
	}

	writeCSV(river.Name+"_"+parm.Query()["data"][0]+".csv", river.WaterLevels)
}

// 河名＋局名
func getRiverName(doc *goquery.Document) string {
	var Name = ""
	var SiteName = ""
	doc.Find("table.horizontal.first tr").Each(func(i int, s *goquery.Selection) {
		title, value := getPairOfColumn(s)
		//fmt.Printf("%d:%s\n", i, value)
		if title == "河川名" {
			Name = value
		}
		if title == "局名" {
			SiteName = value
		}
	})
	fmt.Println("[" + Name + "]")
	return strings.Trim(Name+"-"+SiteName, "\t")
}

func getPairOfColumn(tr *goquery.Selection) (string, string) {
	th := tr.Find("th")
	return th.First().Text(), th.Next().Text()
}

const format = "2006.01.02 15:04"

func getWaterLevels(doc *goquery.Document) []WaterLevel {
	wls := []WaterLevel{}
	doc.Find("table.horizontal.second tr").Each(func(i int, s *goquery.Selection) {
		title, value := getPairOfColumn(s)
		// fmt.Printf("%s:%s\n", title, value)
		t, _ := time.Parse(format, title)
		l, _ := strconv.ParseFloat(value, 64)
		wl := WaterLevel{
			Time:  t,
			Level: l,
		}
		wls = append(wls, wl)
	})
	return wls

}

func buildURL(u string) string {
	return fmt.Sprintf("http://www.bousai-kagawa.jp/%s", u)
}

func writeCSV(filename string, data []WaterLevel) {
	file, err := os.Create("csv/" + filename)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	for _, v := range data {
		ls := []string{v.Time.Format("2006-01-02 15:04:05"), fmt.Sprintf("%g", v.Level)}
		writer.Write(ls)

	}
	writer.Flush()
	fmt.Println("output: csv/" + filename)
}
