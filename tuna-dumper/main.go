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
	"github.com/ismdeep/ismdeep-go-utils/args_util"
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

func ShowHelpMsg() {
	fmt.Println("Usage: tuna-dumper   -url <URL>    -out <output>    [--append]")
}

func main() {
	if args_util.Exists("--help") {
		ShowHelpMsg()
		return
	}

	if !args_util.Exists("-url") || !args_util.Exists("-out") {
		ShowHelpMsg()
		os.Exit(1)
	}

	repoUrl := args_util.GetValue("-url")
	appendFlag := args_util.Exists("--append")

	fmt.Printf("URL: [%s]\n", repoUrl)
	fmt.Printf("Append: %v\n", appendFlag)

	urls := getSubUrls(repoUrl)
	sort.Strings(urls)

	for _, url := range urls {
		fmt.Println(url)
	}
}
