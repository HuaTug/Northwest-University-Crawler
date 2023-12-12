package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func fetch(url string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", "_ga=GA1.2.975660501.1687059913; Hm_lvt_866c9be12d4a814454792b1fd0fed295=1687059912,1687253977; clostip=0")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Http get err:", err)
		return ""
	}
	if resp.StatusCode != 200 {
		fmt.Println("Http status code:", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error", err)
		return ""
	}
	return string(body)
}

func Parse(html string) {
	//替换掉空格 这段代码的意思是将字符串变量 html 中的换行符 (\n) 替换为空字符串 ("")。具体来说，strings.Replace 函数用于在一个字符串中查找指定的子字符串，然后将其替换为新的字符串。在这个例子中，"\n" 是要查找和替换的子字符串，"" 是用于替换的新字符串，-1 表示替换所有匹配的子字符串。
	html = strings.Replace(html, "\n", "", -1)
	// 边栏内容块正则
	re_sidebar := regexp.MustCompile(`<aside id="sidebar" role="navigation">(.*?)</aside>`) //(.*?)表示为任意内容
	//找到边栏内容块
	sidebar := re_sidebar.FindString(html) //html 是一个字符串，代表整个 HTML 文档  使用 re_sidebar.FindString 方法在 html 字符串中查找与正则表达式模式匹配的第一个子字符串，并返回结果作为一个字符串
	//链接正则
	//ToDo :使用 regexp.MustCompile 函数创建了一个正则表达式对象 re_link，该正则表达式的模式是 href="(.*?)"
	re_link := regexp.MustCompile(`href="(.*?)"`)
	//ToDo :使用 re_link.FindAllString 方法在 sidebar 字符串中查找所有与正则表达式模式匹配的子字符串，并返回结果作为一个字符串切片。第二个参数 -1 表示匹配的次数没有限制，即找到所有匹配的子字符串。
	//找到所有链接
	links := re_link.FindAllString(sidebar, -1)
	base_url := "https://gorm.io/zh_CN/docs/"
	for _, v := range links {
		fmt.Printf("url:%v\n", v)
		s := v[6 : len(v)-1]
		url := base_url + s
		fmt.Printf("url:%v\n", url)
		body := fetch(url)
		go Parse2(body)
	}
}

func Parse2(body string) {
	//在HTML页面中进行匹配Content
	body = strings.Replace(body, "\n", "", -1)
	re_content := regexp.MustCompile(`<div class="article">(.*?)</div>`)
	content := re_content.FindString(body)
	//FindString(s) 将返回第一个匹配的数字子字符串。

	//在Content中继续匹配标题Title
	re_title := regexp.MustCompile(`<h1 class="article-title" itemprop="name">(.*?)</h1>`) //ToDo :用于标识获取标题
	title := re_title.FindString(content)
	fmt.Printf("Title: %v\n", title)

	title = title[42 : len(title)-5]
	fmt.Printf("Title: %v\n", title)
	//save(title, content)
	err := saveToDB(title, content)
	if err != nil {
		return
	}
}

func save(title string, content string) {
	//ToDo :"./Log/" 表示为在当前目录下的Log目录中存放到本地文件
	err := os.WriteFile("./Log/"+title+".html", []byte(content), 0644)
	if err != nil {
		panic(err)
	}
}

type GormPage struct {
	gorm.Model
	Title   string `gorm:"title"`
	Content string `gorm:"content"`
}

func saveToDB(title, content string) error {
	dsn := "root:123456@tcp(127.0.0.1:3306)/bug?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	_ = db.AutoMigrate(&GormPage{})
	page := GormPage{
		Title:   title,
		Content: content,
	}
	err = db.Create(&page).Error
	if err == nil {
		return err
	}
	return nil
}
func main() {

	url := "https://gorm.io/zh_CN/docs/"
	s := fetch(url)
	//fmt.Printf("Results: %v\n", s)
	Parse(s)
}
