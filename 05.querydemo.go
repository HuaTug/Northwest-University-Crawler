package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://gorm.io/zh_CN/docs/"
	doc, _ := goquery.NewDocument(url)

	doc.Find(".sidebar-link").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		fmt.Printf("href: %v\n", href)
		detail_url := url + href
		doc, _ = goquery.NewDocument(detail_url)
		title := doc.Find(".article-title").Text()
		fmt.Printf("title: %v\n", title)

		//content, _ := doc.Find(".article").Html()
		//fmt.Printf("Content: %v/n", content)
	})
}
