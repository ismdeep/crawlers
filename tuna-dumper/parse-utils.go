package main

import (
	"errors"
	"github.com/ismdeep/ismdeep-go-utils/args_util"
	"os"
)

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
