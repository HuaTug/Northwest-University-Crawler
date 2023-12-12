package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Inf struct {
	gorm.Model
	Title string `json:"title" gorm:"column:title"`
	Info  string `json:"info" gorm:"column:info"`
}

var (
	db   *gorm.DB
	Infs []*Inf
)

func fetch2(url string) *goquery.Document {
	fmt.Println("Fetch Url", url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Http get err:", err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("Http status code:", resp.StatusCode)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func parseUrls(url string) {
	doc, _ := goquery.NewDocument(url)
	doc.Find(".erji-content-div li").Each(func(index int, ele *goquery.Selection) {
		movieUrl, _ := ele.Find("a").Attr("href")

		title, _ := ele.Find("a").Attr("title")
		movieUrl = "https://jwc.nwu.edu.cn/" + movieUrl[2:]
		fmt.Println(movieUrl)
		fmt.Println(title)
		inf := &Inf{
			Title: title,
			Info:  movieUrl,
		}
		Infs = append(Infs, inf)
	})

}
func init() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/bug?charset=utf8mb4&parseTime=True"
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Print("数据库连接成功!")
	db = d
	db.AutoMigrate(&Inf{})
}

func main() {
	start := time.Now()
	for i := 1; i <= 28; i++ {
		parseUrls("https://jwc.nwu.edu.cn/tzgg1" + "/" + strconv.Itoa(29-i) + ".htm")
	}
	if err := db.Create(Infs).Error; err != nil {
		log.Print("插入数据失败", err)
	}
	elapsed := time.Since(start)
	fmt.Printf("Took %s", elapsed)
}
