package main

/**
 * Usage: tuna-dumper   -url <URL>   -a <leave_tag> -x <remove_tag> -out <output> [--append]
 *
 * Author: L. JIANG <l.jiang.1024@gmail.com>
 *
 */

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/ismdeep/ismdeep-go-utils/args_util"
	"golang.org/x/net/html"
	"log"
	"os"
	"sort"
	"strings"
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
	fmt.Println("Usage: tuna-dumper   -url <URL>    -a <leave_tag>    -x <remove_tag>    -out <output>    [--append]")
}

func parseOutputPath() (string, error) {
	for i := 1; i < len(os.Args)-1; i++ {
		if os.Args[i] == "-out" {
			return os.Args[i+1], nil
		}
	}
	return "", errors.New("[ERROR] --out argument is required")
}

func parseFilterList() []string {
	filterList := make([]string, 0)
	for i := 1; i < len(os.Args)-1; i++ {
		if os.Args[i] == "-a" {
			filterList = append(filterList, os.Args[i+1])
		}
	}
	return filterList
}

func parseFilterRemoveList() []string {
	filterRemoveList := make([]string, 0)
	for i := 1; i < len(os.Args)-1; i++ {
		if os.Args[i] == "-x" {
			filterRemoveList = append(filterRemoveList, os.Args[i+1])
		}
	}
	return filterRemoveList
}

func parseAppendFlag() bool {
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--append" {
			return true
		}
	}
	return false
}

func parseUrl() (string, error) {
	if args_util.Exists("-url") {
		return args_util.GetValue("-url"), nil
	}
	return "", errors.New("[ERROR] -url not found")
}

func filter(url string, filterList []string, filterRemoveList []string) bool {
	for _, filterTag := range filterList {
		if !strings.Contains(url, filterTag) {
			return false
		}
	}

	for _, filterRemoveTag := range filterRemoveList {
		if strings.Contains(url, filterRemoveTag) {
			return false
		}
	}
	return true
}

func main() {
	if args_util.Exists("--help") {
		ShowHelpMsg()
		return
	}

	repoUrl, err := parseUrl()
	if err != nil {
		ShowHelpMsg()
		log.Fatalln(err)
		return
	}

	appendFlag := parseAppendFlag()
	outputPath, err := parseOutputPath()
	if err != nil {
		ShowHelpMsg()
		log.Fatalln(err)
		return
	}

	filterList := parseFilterList()
	filterRemoveList := parseFilterRemoveList()

	log.Printf("URL: [%s]\n", repoUrl)
	log.Printf("Append: %v\n", appendFlag)
	log.Printf("Out: %v\n", outputPath)
	log.Printf("Filter List: %v\n", filterList)
	log.Printf("Filter Remove List: %v\n", filterRemoveList)

	urls := getSubUrls(repoUrl)
	sort.Strings(urls)

	var f *os.File
	if appendFlag {
		f, err = os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("[ERROR]: can not open file.\n")
			return
		}
	} else {
		f, err = os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("[ERROR]: can not open file.\n")
			return
		}
	}

	w := bufio.NewWriter(f)

	for _, url := range urls {
		if filter(url, filterList, filterRemoveList) {
			log.Println(url)
			w.WriteString(fmt.Sprintf("%v\n", url))
		}
	}

	w.Flush()
	f.Close()
}
