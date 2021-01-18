package magnet_source

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/ismdeep/crawlers/magnet-search/structs"
	"golang.org/x/net/html"
	"os"
	"strconv"
	"strings"
	"sync"
)

func YHG66a_Detail(url string) (structs.MagnetInfo, error) {
	res := structs.MagnetInfo{
		Title:       "",
		MagnetUrl:   "",
		CreateTime:  "",
		FileSize:    "",
		DownloadHot: 0,
	}

	var doc *html.Node
	var err error
	doc, err = htmlquery.LoadURL(url)
	if err != nil {
		fmt.Println(url)
		fmt.Fprintln(os.Stderr, err)
		return res, err
	}

	res.Title = htmlquery.Find(doc, `/html/body/div[@id='wrapper']/div[@id='content']/div[@id='wall']/h2`)[0].FirstChild.Data
	res.CreateTime = htmlquery.Find(doc, `//div[@class='fileDetail']//table[@class='detail-table detail-width']/tbody/tr[2]/td[2]`)[0].FirstChild.Data
	res.DownloadHot, _ = strconv.Atoi(htmlquery.Find(doc, `//div[@class='fileDetail']//table[@class='detail-table detail-width']/tbody/tr[2]/td[3]`)[0].FirstChild.Data)
	res.FileSize = htmlquery.Find(doc, `//div[@class='fileDetail']//table[@class='detail-table detail-width']/tbody/tr[2]/td[4]`)[0].FirstChild.Data
	tmp := htmlquery.Find(doc, `//meta[@name="description"]/@content`)[0].FirstChild.Data
	res.MagnetUrl = strings.TrimSpace(tmp[strings.Index(tmp, "magnet"):])

	return res, nil
}

func YHG66a_Search(keyword string) []structs.MagnetInfo {
	var wg sync.WaitGroup

	results := make([]structs.MagnetInfo, 0)

	for i := 1; i <= 3; i++ {
		var doc *html.Node
		var err error
		doc, err = htmlquery.LoadURL(fmt.Sprintf("https://www.yhg66a.xyz/search/%v-%v-time.html", keyword, i))
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return make([]structs.MagnetInfo, 0)
		}

		nodes := htmlquery.Find(doc, `//div[@id='content']//div[@id='wall']//div[@class='search-item detail-width']//h3/a/@href`)
		for _, node := range nodes {
			uri := node.FirstChild.Data
			wg.Add(1)
			go func(__uri__ string) {
				done := false
				cnt := 3
				for !done && cnt > 0 {
					cnt--
					detail, err := YHG66a_Detail(fmt.Sprintf("https://www.yhg66a.xyz%v", uri))
					if err == nil {
						results = append(results, detail)
						done = true
					}
				}
				wg.Done()
			}(uri)
		}
	}

	wg.Wait()

	return results
}
