package anyflip

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func DownloadConfigJSFile(bookURL *url.URL) (string, error) {
	configjsURL, err := url.Parse("https://online.anyflip.com")
	if err != nil {
		return "", err
	}
	configjsURL.Path = path.Join(bookURL.Path, "mobile", "javascript", "config.js")
	resp, err := http.Get(configjsURL.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("received non-200 response:" + resp.Status)
	}
	configjs, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(configjs), nil
}

func GetPageCount(configjs string) (int, error) {
	r := regexp.MustCompile("\"?(bookConfig\\.)?totalPageCount\"?[=:]\"?\\d+\"?")
	match := r.FindString(configjs)

	if strings.Contains(match, "=") {
		match = strings.Split(match, "=")[1]
	} else if strings.Contains(match, ":") {
		match = strings.Split(match, ":")[1]
	} else {
		return 0, errors.New("could not find page count")
	}

	match = strings.ReplaceAll(match, "\"", "")
	return strconv.Atoi(match)
}

func GetBookTitle(configjs string) (string, error) {
	r := regexp.MustCompile(`("?(bookConfig.)?bookTitle"?[=]"(.*?)")|"title":"(.*?)"`)
	match := r.FindString(configjs)

	if strings.Contains(match, "=") {
		match = strings.Split(match, "=")[1]
	} else if strings.Contains(match, ":") {
		match = strings.Split(match, ":")[1]
	} else {
		return "", errors.New("could not find book title")
	}
	match = strings.ReplaceAll(match, "\"", "")

	return match, nil
}