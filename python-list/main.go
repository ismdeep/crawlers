package main

import (
	"bufio"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"log"
	"os"
	"sort"
	"sync"
)

func getDownloadList(url string) []string {
	var urls []string

	var doc *html.Node
	var err error
	for {
		doc, err = htmlquery.LoadURL(url)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		break
	}

	nodes := htmlquery.Find(doc, `//div[@id='touchnav-wrapper']/div[@id='content']/div[@class='container']/section[@class='main-content ']/article[@class='text']/table/tbody/tr/td/a/@href`)
	for _, node := range nodes {
		urls = append(urls, node.FirstChild.Data)
	}

	return urls
}

func main() {
	homeUrl := "https://www.python.org/downloads/"
	var doc *html.Node
	var err error
	for {
		doc, err = htmlquery.LoadURL(homeUrl)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		break
	}

	nodes := htmlquery.Find(doc, `//div[@class='row download-list-widget']/ol[@class='list-row-container menu']/li/span[@class="release-number"]/a/@href`)
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var downloadUrls []string
	for _, node := range nodes {
		wg.Add(1)
		go func(subUrl string) {
			urls := getDownloadList(subUrl)
			mutex.Lock()
			for _, url := range urls {
				downloadUrls = append(downloadUrls, url)
			}
			mutex.Unlock()
			wg.Done()
		}(homeUrl + node.FirstChild.Data)
	}
	wg.Wait()

	sort.Strings(downloadUrls)

	// 寫入文件
	openMode := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	f, err := os.OpenFile("python.txt", openMode, 0666)
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(f)
	for _, downloadUrl := range downloadUrls {
		w.WriteString(downloadUrl + "\n")
	}
	w.Flush()
	f.Close()
}
