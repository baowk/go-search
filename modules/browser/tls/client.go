package tls

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	fhttp "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

func NewTlsClient(reqUrl, proxyUrl string, cookies []*fhttp.Cookie, clientProfile profiles.ClientProfile) (tls_client.HttpClient, error) {
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
	respRedirect, err := fhttp.Get(redirectFullURL)
	if err != nil {
		slog.Error("发起重定向请求失败", "url", redirectFullURL)
		return
	}
	defer respRedirect.Body.Close()
	fmt.Printf("重定向后响应状态码：%d\n", respRedirect.StatusCode) // 通常为 200 OK
}
