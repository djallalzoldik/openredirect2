package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var urls []string
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			urls = append(urls, line)
		}
	}

	var results []map[string]string
	for _, urlStr := range urls {
		urlObj, err := url.Parse(urlStr)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			continue
		}

		req, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}
		req.Header.Set("User-Agent", "Mozilla/5.0")

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			continue
		}

		defer resp.Body.Close()

		location := resp.Header.Get("Location")
		if location != "" {
			vulnerable := false
			for _, param := range urlObj.Query() {
				for _, value := range param {
					if strings.Contains(location, value) {
						vulnerable = true
						break
					}
				}
				if vulnerable {
					break
				}
			}

			if vulnerable {
				result := map[string]string{
					"url":        urlStr,
					"vulnerable": location,
				}
				results = append(results, result)
			}
		}
	}

	jsonBytes, err := json.Marshal(results)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
	fmt.Println(string(jsonBytes))
}
