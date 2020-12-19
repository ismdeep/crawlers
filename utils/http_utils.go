package utils

import (
	"io"
	"io/ioutil"
	"net/http"
)

func GetHttpBody(url string) (io.ReadCloser, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

func GetHttpContent(url string) (string, error) {
	body, err := GetHttpBody(url)
	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
