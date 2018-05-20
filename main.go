package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var precessed = []string{}
var shipNames = map[string]int{}

func isExist(url string) bool {
	for _, v := range precessed {
		if url == v {
			return true
		}
	}
	return false
}

func main() {
	os.MkdirAll("qchar_images", os.ModePerm)

	c := colly.NewCollector()

	// 取得每個陣營的所有船艦
	for i := 1; i <= 7; i++ {
		c.OnHTML("#con_"+strconv.Itoa(i)+" td[style=\"height:80px\"] a[href]", func(e *colly.HTMLElement) {
			url := e.Attr("href")
			if isExist(url) {
				return
			}
			precessed = append(precessed, url)
			e.Request.Visit(url)
		})
	}

	// 取得船艦的 q 圖
	c.OnHTML(".qchar-container img[src]", func(e *colly.HTMLElement) {
		src := e.Attr("src")
		if isExist(src) {
			return
		}
		precessed = append(precessed, src)

		escapeNames := strings.Split(e.Request.URL.String(), "/")
		escapeName := escapeNames[len(escapeNames)-1]
		realName, err := url.QueryUnescape(escapeName)

		if err != nil {
			return
		}

		filename := realName + "-" + strconv.Itoa(shipNames[realName]) + ".png"
		shipNames[realName] += 1
		fmt.Println("Processing..." + filename)

		out, err := os.Create("qchar_images/" + filename)
		defer out.Close()
		if err != nil {
			fmt.Println("Create File Fail: " + err.Error())
			return
		}

		resp, err := http.Get(src)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println("Get File Fail: " + err.Error())
			return
		}

		_, err = io.Copy(out, resp.Body)
	})

	c.Visit("http://wiki.joyme.com/blhx/%E9%98%B5%E8%90%A5")
}
