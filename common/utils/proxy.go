package utils

import (
	"net/http"
	"net/url"
	"time"
)

// ProxyCheckV2 检测proxy是否可用，如：http://username:password@ip:port
func ProxyCheckV2(proxies []string) map[string]bool {
	res := make(map[string]bool, 0)
	for _, p := range proxies {
		// Parse the proxy URL
		proxyURL, err := url.Parse(p)
		if err != nil {
			res[proxyURL.Hostname()] = false
			continue
		}

		// Create HTTP client with proxy
		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
			Timeout: 10 * time.Second,
		}

		// Make request to ipinfo.io
		resp, err := client.Get("http://myip.ipipv.com")
		if err != nil {
			res[proxyURL.Hostname()] = false
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			res[proxyURL.Hostname()] = false
			continue
		}

		res[proxyURL.Hostname()] = true
	}
	return res
}
