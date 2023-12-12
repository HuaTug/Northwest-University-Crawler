package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func HttpGet(url string) (result string, err error) {
	res, err1 := http.Get(url)
	if err1 != nil {
		err = nil //将封装函数内部的错误，传出给调用者
		return
	}
	defer res.Body.Close()

	//循环读取网页数据
	buf := make([]byte, 4096)
	for {
		n, err2 := res.Body.Read(buf)
		//ToDo res.Body 应该是一个表示HTTP响应体的对象。Read(buf) 是一个方法调用，它会尝试从 res.Body 中读取数据，并将读取的数据存储到 buf 中
		if n == 0 {
			fmt.Println("读取网页完成")
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}
		result += string(buf[:n])
	}
	return
}

//ToDo :并发问题：在Working函数中，使用go SpiderPage(i)启动多个goroutine进行爬取。如果在main函数中没有进行等待操作，那么主goroutine可能会在子goroutine还没有完成爬取操作时就退出，导致爬取任务未完成。你可以使用sync.WaitGroup来等待所有爬取任务完成。

func SpiderPage(i int, wg *sync.WaitGroup) {
	defer wg.Done()
	url := "https://jwc.nwu.edu.cn/tzgg1/" + strconv.Itoa(i) + ".htm"
	results, err := HttpGet(url)
	if err != nil {
		fmt.Println("HttpGet err: ", err)
		return
	}
	//fmt.Println("Results=", results)
	//将读到的整网页数据，保存成一个文件
	f, err := os.Create("第" + strconv.Itoa(i) + "页" + ".html")
	if err != nil {
		fmt.Println("Create err:", err)
		return
	}
	f.WriteString(results)
	f.Close() //保存好一个文件，关闭一个文件
}
func Working(start, end int) {
	var wg sync.WaitGroup
	wg.Add(end - start + 1)
	//循环爬取每一页的数据
	for i := start; i <= end; i++ {
		go SpiderPage(i, &wg)
	}
	wg.Wait()
}
func main() {
	var start, end int
	fmt.Print("请输入爬取的起始页:->")
	fmt.Scan(&start)
	fmt.Print("请输入爬取的终止页:->")
	fmt.Scan(&end)
	Working(start, end)
}
