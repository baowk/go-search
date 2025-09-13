package models

import (
	"net/http"
	"strings"
	"time"
)

var headers = map[string]string{
	"accept":                      "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	"accept-language":             "zh-CN,zh;q=0.9",
	"downlink":                    "6",
	"priority":                    "u=0, i",
	"referer":                     "https://www.google.com/",
	"rtt":                         "100",
	"sec-ch-prefers-color-scheme": "light",
	"sec-ch-ua":                   "\"Not;A=Brand\";v=\"99\", \"Google Chrome\";v=\"139\", \"Chromium\";v=\"139\"",
	"sec-ch-ua-arch":              "\"x86\"",
	"sec-ch-ua-bitness":           "\"64\"",
	"sec-ch-ua-form-factors":      "\"Desktop\"",
	"sec-ch-ua-full-version":      "\"139.0.7258.139\"",
	"sec-ch-ua-full-version-list": "\"Not;A=Brand\";v=\"99.0.0.0\", \"Google Chrome\";v=\"139.0.7258.139\", \"Chromium\";v=\"139.0.7258.139\"",
	"sec-ch-ua-mobile":            "?0",
	"sec-ch-ua-model":             "\"\"",
	"sec-ch-ua-platform":          "\"Windows\"",
	"sec-ch-ua-platform-version":  "\"10.0.0\"",
	"sec-ch-ua-wow64":             "?0",
	"sec-fetch-dest":              "document",
	"sec-fetch-mode":              "navigate",
	"sec-fetch-site":              "same-origin",
	"upgrade-insecure-requests":   "1",
	"user-agent":                  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36",
	"x-browser-channel":           "stable",
	"x-browser-copyright":         "Copyright 2025 Google LLC. All rights reserved.",
	"x-browser-year":              "2025",
	"accept-encoding":             "gzip, deflate, br", // Add gzip support
}

type HeaderCache struct {
	Cookie       string            `json:"cookie"`
	UsedNum      int               `json:"usedNum"`
	LastUsedTime time.Time         `json:"lastUsedTime"`
	UserAgent    string            `json:"userAgent"`
	HeaderMap    map[string]string `json:"headerMap"`
}

func (c *HeaderCache) SetRequest(req *http.Request) {
	//初始化默认值
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if c.HeaderMap != nil {
		for k, v := range c.HeaderMap {
			req.Header.Set(k, v)
		}
	}

	if c.UserAgent != "" {
		req.Header.Set("user-agent", c.UserAgent)
	}

	if c.Cookie != "" {
		cookies := strings.Split(c.Cookie, ";")
		for _, cookie := range cookies {
			cookie = strings.TrimSpace(cookie)
			parts := strings.SplitN(cookie, "=", 2)
			if len(parts) == 2 {
				//	fmt.Println(parts[0], parts[1])
				req.AddCookie(&http.Cookie{
					Name:  strings.TrimSpace(parts[0]),
					Value: strings.TrimSpace(parts[1]),
				})
			}
		}
	}
}
