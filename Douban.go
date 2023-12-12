package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type Movie struct {
	gorm.Model
	Title        string  `json:"title" gorm:"column:title"`
	PublishData  string  `json:"publish_data" gorm:"column:publish_data"`
	Score        float64 `json:"Score" gorm:"column:score"`
	CommentCount int64   `json:"comment_count" gorm:"column:comment_count"`
	Quote        string  `json:"quote" gorm:"column:quote"`
}

func (m *Movie) TableName() string {
	return "movie"
}
func NewMovie(title, publishdata, quote string, score float64, commentcount int64) *Movie {
	return &Movie{
		Title:        title,
		PublishData:  publishdata,
		Score:        score,
		CommentCount: commentcount,
		Quote:        quote,
	}
}

var (
	db     *gorm.DB
	movies []*Movie
)

func ClearPlain(str string) string {
	reg := regexp.MustCompile(`\s`)
	//ToDo :正则表达式模式 \s 匹配任何空白字符，包括空格、制表符、换行符等
	return reg.ReplaceAllString(str, "")
}

func GetNumber(str string) string {
	reg := regexp.MustCompile(`\d+`)
	//ToDo 正则表达式模式 \d+ 匹配一个或多个数字字符。\d 是一个预定义的字符类，表示数字字符。+ 是一个量词，表示前面的元素可以出现一次或多次。
	return reg.FindString(str)
}
func init() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/bug?charset=utf8mb4&parseTime=True"
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Print("数据库连接成功!")
	db = d
	db.AutoMigrate(&Movie{})
}
func Run(method, url string, body io.Reader, client *http.Client) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("host", "movie.douban.com")
	resq, err := client.Do(req)
	if err != nil {
		log.Println("发送请求失败")
		return
	}
	if resq.StatusCode != http.StatusOK {
		log.Printf("请求失败，状态码: %d", resq.StatusCode)
		return
	}
	defer resq.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resq.Body)
	if err != nil {
		log.Println("生成goQuery对象失败")
		return
	}
	doc.Find("ol.grid_view li").Each(func(i int, selection *goquery.Selection) {
		title := ClearPlain(selection.Find("span.title").Text())
		year := ClearPlain(GetNumber(selection.Find("div.bd>p").Text()))
		commentCountStr := ClearPlain(GetNumber(selection.Find(".start>span").Eq(3).Text()))
		scoreStr := ClearPlain(selection.Find("span.rating_num").Text())
		quote := ClearPlain(selection.Find(".inq").Text())
		commentCount, _ := strconv.ParseInt(commentCountStr, 10, 64)
		score, _ := strconv.ParseFloat(scoreStr, 64) // 评分可能是小数，所以这里用的是ParseFloat方法
		fmt.Println(title)
		fmt.Println(year)
		fmt.Println(commentCount)
		fmt.Println(score)
		fmt.Println(quote)
		movies = append(movies, NewMovie(title, year, quote, score, commentCount))
		fmt.Println("--------------------------")
	})
}
func main() {
	client := &http.Client{}
	url := "https://movie.douban.com/top250?start=%d&filter="
	method := "GET"
	for i := 1; i <= 10; i++ {
		Run(method, fmt.Sprintf(url, i*25), nil, client)
		time.Sleep(time.Second * 2)
	}
	if err := db.Create(movies).Error; err != nil {
		log.Println("插入数据失败", err.Error())
		return
	}
	log.Println("插入成功")
}
