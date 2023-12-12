package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func fech(url string) string {
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
func main() {
	url := "https://ist.nwu.edu.cn/tzgg.htm"
	s := fech(url)
	fmt.Printf("Results: %v\n", s)
}
