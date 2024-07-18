package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"
)

var (
	baseUri       *url.URL
	urlExpression *regexp.Regexp
	randomizer    *rand.Rand
)

func main() {
	arguments := os.Args[1:]
	if len(arguments) != 1 {
		fmt.Fprintln(os.Stderr, "usage\n\t./task link")
		os.Exit(1)
	}

	var err error
	baseUri, err = url.Parse(arguments[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "usage\n\t./task link")
		os.Exit(1)
	}

	urlExpression = regexp.MustCompile(`(?i)<a[^>]+href=['"]([^'"]+)['"]`)

	client := http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    100,
			IdleConnTimeout: 90 * time.Second,
			MaxConnsPerHost: 10,
		},
		Timeout: 5 * time.Second,
	}

	salt := rand.NewSource(time.Now().Unix())
	randomizer = rand.New(salt)

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	go downloadPage(baseUri.String(), client, waitGroup, 1)

	waitGroup.Wait()
}

func downloadPage(uri string, client http.Client, wg *sync.WaitGroup, depth int) {
	defer wg.Done()
	if depth >= 3 {
		fmt.Fprintln(os.Stdout, "reached the end")
		return
	}

	request, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while creating request: %v", err)
		return
	}

	response, err := client.Do(request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while doing request: %v", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		fmt.Fprintf(os.Stderr, "error: status code %d", response.StatusCode)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while reading response body: %v", err)
		return
	}

	file, err := os.Create(fmt.Sprintf("%d.txt", randomizer.Intn(1000)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while creating save file: %v", err)
		return
	}
	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while copying to save file: %v", err)
		return
	}

	urls, err := findUrls(body)
	if err != nil {
		return
	}
	ticker := time.NewTicker(6 * time.Second)
	fmt.Println(urls)
	for idx, u := range urls {
		if idx > 5 {
			break
		}
		u = strings.TrimSpace(u)

		if u == uri {
			continue
		}
		<-ticker.C
		wg.Add(1)
		go downloadPage(baseUri.String()+u, client, wg, depth+1)
	}
}

func findUrls(data []byte) ([]string, error) {
	paths := urlExpression.FindAllStringSubmatch(string(data), -1)
	urls := make([]string, len(paths))

	for idx, path := range paths {
		pathString := path[1]
		if checkPath(pathString) || slices.Contains(urls, pathString) {
			continue
		}

		urls[idx] = pathString
	}

	return urls, nil
}

func checkPath(path string) bool {
	return strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "http://") || path == ""
}
