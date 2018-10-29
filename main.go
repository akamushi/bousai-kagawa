package main

import (
	"log"
	"net/http"

	"github.com/akamushi/bousai-kagawa/river"

	"github.com/PuerkitoBio/goquery"
)

// ここからスタート
func main() {
	// 水位マップページにアクセス
	res, err := http.Get("http://www.bousai-kagawa.jp/suii-map.html")
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

	findRiverPlot(doc)
}

func findRiverPlot(doc *goquery.Document) {
	// ページ内の河川プロットをすべてに対してURLを取得し、
	// ダウンロードを実施する。
	doc.Find("div.info.cf div.plot").Each(func(i int, s *goquery.Selection) {
		u, _ := s.Find("a").Attr("href")

		// fmt.Println(url)
		river.Download(u)

	})
}
