package handler

import (
	"bytes"
	"compress/gzip"
	"dilu/common/utils"
	"dilu/modules/browser"
	"dilu/modules/search/google/enums"
	"dilu/modules/search/service/dto"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/brotli"
)

func ToSearch(params *dto.SearchReq, res *dto.SearchResp) (int, error) {
	begin := time.Now()
	if params.Engine == "" {
		params.Engine = enums.EngineGoogle
	}
	if params.Num <= 0 {
		params.Num = 10
	} else if params.Num > 100 {
		params.Num = 100
	}

	no := utils.GenString()
	res.SearchMetadata.ID = no
	res.SearchMetadata.CreatedAt = begin.Format("2006-01-02 15:04:05")

	res.URL = params.GetGoogleUrl()
	res.SearchMetadata.GoogleURL = res.URL
	res.SearchInformation = append(res.SearchInformation, dto.GetSearchInformationItem{
		QueryDisplayed: params.Q,
	})

	proxyUrl, err := GetProxy()
	if err != nil {
		return 1001, err
	}
	header, err := GetReqHeader()
	if err != nil {
		BackProxy(proxyUrl)
		return 1000, err
	}

	data, sCode, err := SearchV2(params, proxyUrl, header)

	if sCode != enums.Success {
		slog.Error("Search", "sCode", sCode, "proxyUrl", proxyUrl, "err", err)
		switch sCode {
		case enums.ErrParams, enums.ErrDecode, enums.ErrParseHtml:
			BackProxy(proxyUrl)
			BackReqHeader(header)
		case enums.ErrProxy, enums.ErrConn:
			SetFailProxy(proxyUrl)
			BackReqHeader(header)
		case enums.Err429:
			SetFailProxy(proxyUrl)
			BackReqHeader(header)
		case enums.ErrRiskControl, enums.ErrRiskControlClickCode, enums.ErrRiskControlJavaScriptCode:
			SetFailProxy(proxyUrl)
		default:
			SetFailProxy(proxyUrl)
			SetFailReqHeader(header)
		}
	} else {
		BackProxy(proxyUrl)
		BackReqHeader(header)
		if params.Html == "1" {
			res.RawHtml = string(data)
		} else {
			pCode, err := ParseHtmlAll(data, params, res)
			if pCode != enums.Success {
				return pCode, err
			}
		}
	}

	end := time.Now()
	diff := end.Sub(begin)
	res.SearchMetadata.TotalTimeTaken = diff.Seconds()
	res.SearchMetadata.ProcessedAt = end.Format("2006-01-02 15:04:05")
	res.SearchParameters = *params
	return sCode, err
}

func SearchV2(params *dto.SearchReq, proxyUrl string, header *SimpleCookie) ([]byte, int, error) {
	data, code, err := Get(params.GetGoogleUrl(), proxyUrl, header)
	if code != enums.Success {
		slog.Error("Search", "code", code, "proxyUrl", proxyUrl)
		saveToFileWithDir(data, fmt.Sprintf("htmls/%d", code), fmt.Sprintf("err_%s_%d.html", time.Now().Format("20060102"), time.Now().UnixNano()))
	} else {
		if len(data) < 100_000 {
			saveToFileWithDir(data, "htmls/err", fmt.Sprintf("err_%s_%d.html", time.Now().Format("20060102"), time.Now().UnixNano()))
			return nil, enums.ErrRiskControl, errors.New("risk control")
		} else {
			saveToFileWithDir(data, "htmls/ok", fmt.Sprintf("ok_%s_%d.html", time.Now().Format("20060102"), time.Now().UnixNano()))
		}
	}

	return data, code, err
}

func Search(params *dto.SearchReq, proxyUrl string, header *SimpleCookie, device string) ([]byte, int, error) {
	pUrl, err := url.Parse(proxyUrl)
	if err != nil {
		return nil, enums.ErrProxy, err
	}

	req, err := http.NewRequest(http.MethodGet, params.GetGoogleUrl(), nil)
	if err != nil {
		return nil, enums.ErrConn, err
	}

	SetRequest(req, header)
	req.Header.Set("User-Agent", browser.GetUa(device, "chrome", ""))

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(pUrl),
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, enums.ErrConn, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slog.Error("Search", "StatusCode", resp.StatusCode, "proxyUrl", proxyUrl)
		return nil, resp.StatusCode, errors.New("429")
	}
	//t2 := time.Now()

	// Handle different compression formats
	var bodyReader io.Reader
	contentEncoding := resp.Header.Get("Content-Encoding")

	switch {
	case strings.Contains(contentEncoding, "br"):
		// Handle Brotli compression
		bodyReader = brotli.NewReader(resp.Body)
	case strings.Contains(contentEncoding, "gzip"):
		// Handle Gzip compression
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Printf("Failed to create gzip reader, using raw response: %v\n", err)
			bodyReader = resp.Body
		} else {
			defer gzipReader.Close()
			bodyReader = gzipReader
		}
	default:
		// No compression or unsupported compression
		bodyReader = resp.Body
	}

	data, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil, enums.ErrDecode, err
	}

	if len(data) < 100_000 {
		saveToFileWithDir(data, "htmls/err", fmt.Sprintf("err_%s_%d.html", time.Now().Format("20060102"), time.Now().UnixNano()))
		if bytes.Contains(data, []byte(enums.ErrRiskControlClick)) {
			//saveToFileWithDir(data, "errhtml", fmt.Sprintf("ge_click_%s_%d.html", time.Now().Format("20060102"), time.Now().UnixNano()))
			return nil, enums.ErrRiskControlClickCode, fmt.Errorf("click risk control")
		}
		return data, enums.ErrRiskControl, fmt.Errorf("keyword not found in response: %s", params.Q)
	} else {
		saveToFileWithDir(data, "htmls/ok", fmt.Sprintf("%s_%s_%d.html", device, time.Now().Format("20060102"), time.Now().UnixNano()))
	}

	return data, enums.Success, err
}

func SetRequest(req *http.Request, sc *SimpleCookie) {
	if sc != nil && sc.C != "" {
		cookies := strings.Split(sc.C, ";")
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

func saveToFileWithDir(data []byte, dir, filepath string) error {
	// Create directory if it doesn't exist
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(dir + "/" + filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func ParseHtmlAll(html []byte, params *dto.SearchReq, res *dto.SearchResp) (int, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return enums.ErrParseHtml, err
	}

	seleRoot := doc.Find("#center_col")

	alsoAskEl := &goquery.Selection{}

	currentNum := 0

	rso := seleRoot.Find("#rso div.MjjYud")
	if rso.Length() == 1 {
		rso = seleRoot.Find("#kp-wp-tab-overview > div")
	}

	//resultList := make([]*models.KeywordResultItem, 0)

	page := 1
	if params.Start > 0 {
		page = int(math.Ceil(float64(params.Start)/float64(params.Num))) + 1
	}

	rso.Each(func(i int, s *goquery.Selection) {
		//fmt.Println("=======", i, s.Text())
		source := s.Find("div > div > div:nth-child(1) > div > div > span > a > div > div > div > div:nth-child(1) > span").Text()
		// #rso > div:nth-child(3) > div > div > div > div > div > span > a > div > div > div > div.byrV5b > cite
		sourceUrl := s.Find("div > div > div:nth-child(1) > div > div > span > a > div > div > div > div > cite").Text()
		sourceLogo, _ := s.Find("div > div > div:nth-child(1) > div > div > span > a > div > div > span > div > img").Attr("src")
		title := s.Find("div > div > div:nth-child(1) > div > div > span > a > h3").Text()
		titleUrl, _ := s.Find("div > div > div:nth-child(1) > div > div > span > a").Attr("href")
		desc := s.Find("div > div > div:nth-child(2) > div > span").Text()
		d := s.Find("div > div > div > div > div > div > div > span > span:not([class]):not([style]):not([id])").Text()

		if titleUrl != "" && title != "" {
			res.OrganicResults = append(res.OrganicResults, dto.OrganicResultsItem{
				Position:         i + 1,
				Title:            title,
				Url:              titleUrl,
				Description:      strings.TrimRight(desc, " "), // "..."
				SourceLogo:       sourceLogo,
				OriginNavigation: sourceUrl,
				OriginSite:       source,
				Date:             d,
			},
			)
			//resultList = append(resultList, )
			currentNum++
		}
		if page == 1 {
			askHtml, _ := s.Html()
			if strings.Contains(askHtml, `<span>People also ask</span>`) || strings.Contains(askHtml, `<span>相关问题</span>`) {
				alsoAskEl = s
			}
		}
	})
	if page == 1 {
		relatedSearchArr := make([]dto.RelatedSearchesDataItem, 0)
		// 关联搜索  用户还搜索了
		// #bres > div > div > div > div > div > div > div:nth-child(1) > div:nth-child(1) > div > div > a
		seleRoot.Find("#bres div.wyccme").Each(func(i int, s *goquery.Selection) {
			relationUrl, _ := s.Parent().Parent().Attr("href")

			relatedSearchArr = append(relatedSearchArr, dto.RelatedSearchesDataItem{
				Title: s.Text(),
				URL:   handleGoogleUrl(relationUrl, params, res),
			})
		})
		if len(relatedSearchArr) > 0 {
			res.RelatedSearches = dto.RelatedSearchesData{
				Data: relatedSearchArr,
			}
		}

		// 分页数据
		if params.Num <= 0 {
			params.Num = 10
		}
		maxPage := math.Ceil(float64(100) / float64(params.Num))
		pageMap := make(map[string]string)
		pageNextUrl := ""
		seleRoot.Find("#botstuff > div > div > table > tbody > tr > td > a").Each(func(i int, s *goquery.Selection) {
			pageUrlTmp, _ := s.Attr("href")

			pageUrlTmp = handleGoogleUrl(pageUrlTmp, params, res)
			if s.Text() != "" && pageUrlTmp != "" { // 数字分页
				pageNumInt, pageNumIntErr := strconv.Atoi(s.Text())
				if pageNumIntErr == nil {
					if int(maxPage) >= pageNumInt { // 限制最大页数
						if pageNumInt == page+1 {
							pageNextUrl = pageUrlTmp
						}
						pageMap[s.Text()] = pageUrlTmp
					}
				}
			}
		})
		res.Pagination = dto.GetSearchPagination{
			Current:    page,
			Next:       pageNextUrl,
			OtherPages: pageMap,
		}

		//#kp-wp-tab-overview > div > div > div > div > div > div > div > div > div > div > div > div > div > div > div > span > span
		//#rso > 				                                                  div > div > div > div > div > div > div > span > span
		// also ask 相关问题
		if alsoAskEl != nil {
			//slog.Info("OnHtml", "also-ask", alsoAskEl.Find("div > div > div > div > div > div > span > span").Length())
			alsoAskEl.Find("div > div > div > div > div > div > span > span").Each(func(i int, ask *goquery.Selection) {
				//slog.Info(funcName, "also-ask", ask.Text())
				if ask.Text() != "" && (strings.HasSuffix(ask.Text(), "?") || strings.HasSuffix(ask.Text(), "？")) {
					res.AlsoAsk = append(res.AlsoAsk, dto.AlsoAskItem{
						Question: ask.Text(),
					})
				}
			})
		}
	}

	return enums.Success, nil
}

func handleGoogleUrl(url string, params *dto.SearchReq, res *dto.SearchResp) string {
	if params.GoogleDomain != "" {
		if strings.Contains(url, params.GoogleDomain) {
			return url
		}
		return fmt.Sprintf("https://www.%s/%s", params.GoogleDomain, strings.TrimLeft(url, "/"))
	} else {
		if strings.Contains(url, "google.com") {
			return url
		}
		return fmt.Sprintf("https://www.%s/%s", "google.com", strings.TrimLeft(url, "/"))
	}
}
