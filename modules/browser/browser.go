package browser

import (
	"fmt"
	"math/rand/v2"

	fhttp "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"

	"dilu/modules/browser/tls"
	"dilu/modules/browser/ua"
)

var (
	Browsers = map[string][]string{
		"chrome":  {"107", "108", "109", "110", "111", "112", "117", "120", "124", "131", "133"},
		"firefox": {"102", "104", "105", "108", "110", "117", "120", "123", "132", "133", "135"},
		// "safari":  {"15_6_1", "16_0"},
		// "opera":   {"89", "90", "91"},
	}

	//"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	Chrome_UA  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s.0.0.0 Safari/537.36"
	Edge_UA    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s.0.0.0 Safari/537.36 Edg/%s.0.0.0"
	Firefox_UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:%s.0) Gecko/20100101 Firefox/%s.0"
)

func NewClient(reqUrl, proxyUrl string, cookies []*fhttp.Cookie, browser, ver string) (tls_client.HttpClient, string, error) {
	if browser == "" {
		browser = "chrome"
	}
	prof, curUa := GetProfileAndUa(browser, ver)
	c, err := tls.NewTlsClient(reqUrl, proxyUrl, cookies, prof)
	return c, curUa, err
}

func GetProfileAndUa(browser, ver string) (profiles.ClientProfile, string) {
	if ver == "" {
		ver = Browsers[browser][rand.IntN(len(Browsers[browser]))]
	}
	name := fmt.Sprintf("%s_%s", browser, ver)
	prof := profiles.MappedTLSClients[name]
	curUa := GetUa("", browser, ver)
	return prof, curUa
}

func GetUa(device, browser, ver string) string {
	switch device {
	case "android":
		return "Mozilla/5.0 (Linux; Android 8.0; Pixel 2 Build/OPD3.170816.012) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.1639.1179 Mobile Safari/537.36"
	case "iphone":
		return "Mozilla/5.0 (iPhone; CPU iPhone OS 18_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/140.0.7339.122 Mobile/15E148 Safari/604.1"
	case "ipad":
		return "Mozilla/5.0 (iPad; CPU OS 15_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/84.0.4147.71 Mobile/15E148 Safari/604.1"
	case "safari":
		return "Mozilla/5.0 (iPhone; CPU iPhone OS 18_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.5 Mobile/15E148 Safari/604.1"
	}

	if ver == "" {
		ver = Browsers[browser][rand.IntN(len(Browsers[browser]))]
	}
	var curUa string
	switch browser {
	case "chrome":
		curUa = fmt.Sprintf(Chrome_UA, ver)
	case "firefox":
		curUa = ua.GenerateWindowsFirefoxUA("", true, ver, ver)
	case "edge":
		curUa = fmt.Sprintf(Edge_UA, ver)
	default:
		curUa = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
	}
	return curUa
}
