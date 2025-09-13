package handler

import (
	"dilu/modules/search/google/enums"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strings"

	fhttp "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

var (
	Browsers = map[string][]string{
		"chrome":  {"107", "108", "109", "110", "111", "112", "117", "120", "124", "131", "133"},
		"firefox": {"102", "104", "105", "108", "110", "117", "120", "123", "132", "133", "135"},
		"safari":  {"15_6_1", "16_0"},
		"opera":   {"89", "90", "91"},
	}

	//"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	Chrome_UA  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s.0.0.0 Safari/537.36"
	Edge_UA    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s.0.0.0 Safari/537.36 Edg/%s.0.0.0"
	Firefox_UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:%s.0) Gecko/20100101 Firefox/%s.0"
)

func Get(reqUrl, proxyUrl string, sc *SimpleCookie) ([]byte, int, error) {
	prof, ua := GetProfileAndUa("chrome")
	//fmt.Println("profile:", prof, "ua", ua)

	var cookies []*fhttp.Cookie
	if sc.C != "" {
		cookies = ToCookies(sc.C)
	}

	client, err := NewTlsClient(reqUrl, proxyUrl, cookies, prof)
	if err != nil {
		return nil, 500, err
	}
	req, err := fhttp.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, enums.ErrConn, err
	}

	req.Header = fhttp.Header{
		// "accept":          {"*/*"},
		// "accept-language": {"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"},
		"user-agent": {ua},
		// fhttp.HeaderOrderKey: {
		// 	"accept",
		// 	"accept-language",
		// 	"user-agent",
		// },
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, enums.ErrConn, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, errors.New("connection error")
	} else if resp.StatusCode == http.StatusFound {
		slog.Error("302", "statusCode", resp.StatusCode)
		To302(resp, reqUrl)
	}

	cs := resp.Cookies()
	if len(cs) > 0 {
		slog.Info("cookies", "len", len(cs))
		for _, c := range cs {
			for _, c2 := range cookies {
				if c.Name == c2.Name {
					c2.Value = c.Value
				}
			}
		}
		var strC string
		for _, c := range cookies {
			if c.Name == "NID" {
				strC += c.Name + "=" + c.Value + ";"
			}
		}
		sc.C = strC
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, enums.ErrConn, err
	}

	return data, enums.Success, nil
}

func NewTlsClient(reqUrl, proxyUrl string, cookies []*fhttp.Cookie, clientProfile profiles.ClientProfile) (tls_client.HttpClient, error) {
	//fmt.Println("[aaaaaaaaaaaaaaaaa]", reqUrl, proxyUrl, len(cookies), clientProfile)
	rUrl, err := url.Parse(reqUrl)
	if err != nil {
		return nil, err
	}

	jar := tls_client.NewCookieJar()
	if len(cookies) > 0 {
		jar.SetCookies(rUrl, cookies)
	}

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(10),
		tls_client.WithClientProfile(clientProfile),
		//tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
		tls_client.WithProxyUrl(proxyUrl),
		tls_client.WithRandomTLSExtensionOrder(),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		slog.Error("NewTlsClient", "err", err)
	}

	return client, nil
}

func GetProfileAndUa(browser string) (profiles.ClientProfile, string) {
	ver := Browsers[browser][rand.IntN(len(Browsers[browser]))]
	//ver := Browsers[browser][7]

	name := fmt.Sprintf("%s_%s", browser, ver)
	fmt.Println("name", name)
	prof := profiles.MappedTLSClients[name]
	ua := fmt.Sprintf(Chrome_UA, ver)
	return prof, ua
}

func ToCookies(strCookies string) []*fhttp.Cookie {
	if strCookies == "" {
		return nil
	}
	hcs := make([]*fhttp.Cookie, 0)
	cookies := strings.Split(strCookies, ";")
	for _, cookie := range cookies {
		cookie = strings.TrimSpace(cookie)
		parts := strings.SplitN(cookie, "=", 2)
		if len(parts) == 2 {
			//	fmt.Println(parts[0], parts[1])
			hcs = append(hcs, &fhttp.Cookie{
				Name:  strings.TrimSpace(parts[0]),
				Value: strings.TrimSpace(parts[1]),
			})
		}
	}
	return hcs
}

func To302(resp *fhttp.Response, initialURL string) {
	redirectURLStr := resp.Header.Get("Location")
	if redirectURLStr == "" {
		slog.Error("302 响应缺少 Location 头，无法重定向")
		return
	}

	initialParsedURL, err := url.Parse(initialURL)
	if err != nil {
		slog.Error("302 响应缺少 Location 头，无法重定向", "url", initialURL)
		return
	}
	// 拼接相对路径为完整 URL（如 initialURL 是 http://a.com，Location 是 /b → 结果是 http://a.com/b）
	redirectParsedURL, err := initialParsedURL.Parse(redirectURLStr)
	if err != nil {
		slog.Error("302 响应缺少 Location 头，无法重定向", "url", redirectURLStr)
		return
	}
	redirectFullURL := redirectParsedURL.String()

	// 6. 输出 302 信息
	fmt.Printf("检测到 302 临时重定向\n")
	fmt.Printf("初始 URL：%s\n", initialURL)
	fmt.Printf("重定向目标 URL：%s\n", redirectFullURL)

	// （可选）7. 向重定向目标 URL 发起第二次请求
	respRedirect, err := http.Get(redirectFullURL)
	if err != nil {
		slog.Error("发起重定向请求失败", "url", redirectFullURL)
		return
	}
	defer respRedirect.Body.Close()
	fmt.Printf("重定向后响应状态码：%d\n", respRedirect.StatusCode) // 通常为 200 OK
}
