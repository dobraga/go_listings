package main

import "fmt"

func CreateHeaders(origin string) map[string]string {
	return map[string]string{
		"user-agent":       "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
		"accept-language":  "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7",
		"sec-fetch-site":   "same-site",
		"accept":           "application/json",
		"sec-fetch-dest":   "empty",
		"sec-ch-ua-mobile": "?0",
		"sec-fetch-mode":   "cors",
		"origin-ua-mobile": "?0",
		"referer":          fmt.Sprintf("https://www.%s.com.br", origin),
		"origin":           fmt.Sprintf("https://www.%s.com.br", origin),
		"x-domain":         fmt.Sprintf("www.%s.com.br", origin),
	}
}
