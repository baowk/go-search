package handler

import (
	"dilu/modules/browser"
	"dilu/modules/browser/tls"
	"dilu/modules/search/google/enums"
	"errors"
	"io"
	"log/slog"
	"math/rand/v2"

	fhttp "github.com/bogdanfinn/fhttp"
)

func Get(reqUrl, proxyUrl string, sc *SimpleCookie) ([]byte, int, error) {

	var cookies []*fhttp.Cookie
	if sc != nil && sc.C != "" {
		cookies = tls.ToCookies(sc.C)
	}

	b := "chrome"
	r := rand.IntN(10)
	if r%3 == 0 {
		b = "firefox"
	}

	client, ua, err := browser.NewClient(reqUrl, proxyUrl, cookies, b, "")
	if err != nil {
		return nil, 500, err
	}
	req, err := fhttp.NewRequest(fhttp.MethodGet, reqUrl, nil)
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
	if resp.StatusCode != fhttp.StatusOK {
		return nil, resp.StatusCode, errors.New("connection error")
	} else if resp.StatusCode == fhttp.StatusFound {
		slog.Error("302", "statusCode", resp.StatusCode)
		tls.To302(resp, reqUrl)
	}

	cs := resp.Cookies()
	if len(cs) > 0 {
		//slog.Info("cookies", "len", len(cs))
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
