package main

/**
 * Usage: tuna-dumper   <URL>    <RepoName> [--append]
 *
 * Author: L. JIANG <l.jiang.1024@gmail.com>
 *
 */

import (
	"bufio"
	"fmt"
	"github.com/antchfx/htmlquery"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

func getSubUrls(__url__ string) (urls []string) {
	var res *http.Response
	var err error
	for {
		res, err = http.Get(__url__)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	doc, err := htmlquery.Parse(res.Body)
	if err != nil {
		log.Fatal(err)
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
	if len(os.Args) < 3 {
		fmt.Print("Usage: tuna-dumper   <URL>    <RepoName> [--append]")
		os.Exit(1)
	}

	repoUrl := os.Args[1]
	repoName := os.Args[2]
	appendFlag := false
	if len(os.Args) > 3 && os.Args[3] == "--append" {
		appendFlag = true
	}

	urls := getSubUrls(repoUrl)
	sort.Strings(urls)

	fileName := strings.Split(repoName, "/")[0]

	openMode := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	if appendFlag {
		openMode = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	}

	f, err := os.OpenFile(fileName+".txt", openMode, 0666)
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(f)
	for _, url := range urls {
		w.WriteString(url + "\n")
	}
	w.Flush()
	f.Close()
}
