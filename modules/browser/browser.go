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
	curUa := GetUa(browser, ver)
	return prof, curUa
}

func GetUa(browser, ver string) string {
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
