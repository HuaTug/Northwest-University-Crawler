package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

func getDoc1() {
	url := "https://gorm.io/zh_CN/docs/"
	dom, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatalln(err)
	}

	dom.Find(".sidebar-link").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href") //用于取属性
		text := s.Text()          //用于取内容
		fmt.Println(i, href, text)
	})
}

func getDoc2() {
	client := &http.Client{}
	url := "https://gorm.io/zh_CN/docs/"
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	dom, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatalln(err)
	}
	dom.Find("")
}
func getDoc3() {
	html :=
		`<body>
			<div id="div1">XuZh</div>
			<span id="span">SPAN2 </span>"
			<div>DIV2</div>
			<div class="name"> DIV2 </div>
			</body>
			`
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalln(err)
	}
	//ToDo: 这种模式是根据元素名进行查找的操作   元素名称选择器
	dom.Find("div").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
	})
	//ToDo:	根据id进行选择 #{id} id选择器
	dom.Find("#div1").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
	})
	dom.Find("#span").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
	})
	//ToDo: “.{class}"用于获取class类型的数据 class类型选择器
	dom.Find(".name").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
	})
}
func main() {
	getDoc3()
}
