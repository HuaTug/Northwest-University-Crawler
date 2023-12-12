package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"math/rand"
)

type ReplyContainer struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Replies []struct { //评论
			Member struct { //评论用户
				Mid   string `json:"mid"`   //用户id
				Uname string `json:"uname"` //用户姓名
				Sex   string `json:"sex"`   //性别
			} `json:"member"`
			Content struct {
				Message  string        `json:"message"` //评论内容
				Members  []interface{} `json:"members"`
				MaxLine  int           `json:"max_line"`
				Contents struct {
					Message string        `json:"message"`
					Members []interface{} `json:"members"`
					MaxLine int           `json:"max_line"`
				}
			} `json:"content,omitempty"`

			ReplyControl struct {
				MaxLine           int    `json:"max_line"`
				SubReplyEntryText string `json:"sub_reply_entry_text"`
				SubReplyTitleText string `json:"sub_reply_title_text"`
				TimeDesc          string `json:"time_desc"` //评论发布时间
			} `json:"reply_control"`
		} `json:"replies"`
	} `json:"data"`
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandomString 生成一个随机的user-agent
func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func main() {
	c := colly.NewCollector()

	//设置请求头
	c.OnRequest(func(req *colly.Request) {
		req.Headers.Set("authority", "api.bilibili.com")
		req.Headers.Set("accept", "application/json, text/plain, */*")
		req.Headers.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
		req.Headers.Set("cookie", "buvid3=17AE2C8A-B723-FCFD-CA24-BB42886ADC1720800infoc; b_nut=1676006620; i-wanna-go-back=-1; _uuid=AD1102BD9-4377-AEEC-2CDC-10BFF532103E1C20715infoc; DedeUserID=347658083; DedeUserID__ckMd5=c15196b1de6416c6; rpdid=|(JY)RJu)JRJ0J'uY~Y||Ym|); b_ut=5; nostalgia_conf=-1; buvid4=6CEBCDF6-C2BC-E661-A89D-EF5488939A6D21670-023021013-FmUiCC4TrOvqFWiIDh%2F07Q%3D%3D; LIVE_BUVID=AUTO2916763501096807; is-2022-channel=1; hit-dyn-v2=1; blackside_state=0; CURRENT_BLACKGAP=0; SESSDATA=9e31d506%2C1694475764%2Ca17c3%2A32; bili_jct=2318dd41406535428217db54bc617364; CURRENT_PID=982e3db0-c974-11ed-9abb-ab8a696d4338; hit-new-style-dyn=1; buvid_fp_plain=undefined; home_feed_column=4; browser_resolution=1280-649; header_theme_version=CLOSE; FEED_LIVE_VERSION=V8; i-wanna-go-feeds=-1; CURRENT_FNVAL=4048; bili_ticket=eyJhbGciOiJFUzM4NCIsImtpZCI6ImVjMDIiLCJ0eXAiOiJKV1QifQ.eyJleHAiOjE2OTE5ODczMzgsImlhdCI6MTY5MTcyODEzOCwicGx0IjotMX0.IFuNJxp7-eYGGucUiHOGki302PJCOy_DTDnGst8CU22PgItGELCknDv-ttY2W-FcnLZkNlf3lyFSCck9GygWuPLa_sQx5pOXNIdQjUk_pvjyMIwXJvUh7hUBlnhm-dMO; bili_ticket_expires=1691987338; fingerprint=1901b21587cf76da1fcebf71b774fb7a; buvid_fp=30e4f3affecaec8758b47d527c385458; sid=5ng2rgm3; CURRENT_QUALITY=80; b_lsid=7A87A66F_189EDBBC802; bp_video_offset_347658083=829243419043823638; PVID=9")
		req.Headers.Set("origin", "https://www.bilibili.com")
		req.Headers.Set("referer", "https://www.bilibili.com/bangumi/play/ep327584?spm_id_from=333.337.0.0&from_spmid=666.25.episode.0")
		req.Headers.Set("sec-ch-ua", `"Not/A)Brand";v="99", "Microsoft Edge";v="115", "Chromium";v="115"`)
		req.Headers.Set("sec-ch-ua-mobile", "?0")
		req.Headers.Set("sec-ch-ua-platform", `"Windows"`)
		req.Headers.Set("sec-fetch-dest", "empty")
		req.Headers.Set("sec-fetch-mode", "cors")
		req.Headers.Set("sec-fetch-site", "same-site")
		req.Headers.Set("user-agent", RandomString())
	})
	//结构体 用来存放评论数据
	container := ReplyContainer{}
	//c := colly.NewCollector() c是怎么来的
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("response received", r.StatusCode) //打印状态码 成功访问为200
		err := json.Unmarshal(r.Body, &container)      //r.Body为得到的json数据 []byte类型 进行反序列化 字节序列转化成对象
		if err != nil {
			fmt.Println("error", err)
			log.Fatal(err)
		}
	})
	//访问url
	c.Visit("https://api.bilibili.com/x/v2/reply/wbi/main?oid=420981979&type=1&mode=3&pagination_str=%7B%22offset%22:%22%22%7D&plat=1&seek_rpid=&web_location=1315875&w_rid=d0716bedfae00abbbedeccee88bfde56&wts=1702404449")

	for i, reply := range container.Data.Replies {
		fmt.Println(i, "姓名", reply.Member.Uname, "内容", reply.Content.Message)
		fmt.Println(reply.Content.Contents.Message)
		fmt.Println(reply.ReplyControl.TimeDesc, reply.ReplyControl.SubReplyEntryText, reply.ReplyControl.SubReplyEntryText)
	}
}
