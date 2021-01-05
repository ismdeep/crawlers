package main

/**
 * Usage: tuna-dumper   <URL>    <RepoName> [--append]
 *
 * Author: L. JIANG <l.jiang.1024@gmail.com>
 *
 */

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"os"
	"sort"
	"sync"
)

func getSubUrls(__url__ string) (urls []string) {
	var doc *html.Node
	var err error
	for {
		doc, err = htmlquery.LoadURL(__url__)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		break
	}
	nodes := htmlquery.Find(doc, `//tr/td[@class="link"]/a/@href`)

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	for _, node := range nodes {
		subNodeName := node.FirstChild.Data
		if subNodeName == ".." || subNodeName == "../" {
			continue
		}
		if subNodeName[len(subNodeName)-1:] == "/" {
			wg.Add(1)
			go func() {
				resultUrls := getSubUrls(__url__ + subNodeName)
				for _, resultUrl := range resultUrls {
					mutex.Lock()
					urls = append(urls, resultUrl)
					mutex.Unlock()
				}
				wg.Done()
			}()
		} else {
			urls = append(urls, __url__+subNodeName)
		}
	}
	wg.Wait()

	return urls
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: tuna-dumper   <URL>")
		os.Exit(1)
	}

	repoUrl := os.Args[1]

	urls := getSubUrls(repoUrl)
	sort.Strings(urls)

	for _, url := range urls {
		fmt.Println(url)
	}
}
