package main

import (
	"fmt"
	magnet_source "github.com/ismdeep/crawlers/magnet-search/magnet-source"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: magnet-search <keyword>")
		return
	}

	keyword := os.Args[1]

	for _, item := range magnet_source.YHG66a_Search(keyword) {
		fmt.Println(item)
	}
}
