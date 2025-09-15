package apis

import (
	"bufio"
	"dilu/modules/search/service"
	"dilu/modules/search/service/dto"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/baowk/dilu-core/core/base"
	"github.com/gin-gonic/gin"
)

type SearchApi struct {
	base.BaseApi
}

var ApiSearchApi = SearchApi{}

// OpenSearch google search
//
//	@Summary	google search
//	@Tags		open
//	@Accept		application/x-www-form-urlencoded
//	@Product	application/json
//	@Param		api_key			query		string								true	"Parameter api_key open search required"
//	@Param		engine			query		string								false	"Parameter engine Currently only supports google"
//	@Param		q				query		string								true	"Parameter defines the query you want to search. You can use anything that you would use in a regular Google search. e.g. inurl:, site:, intitle:. We also support advanced search query parameters such as as_dt and as_eq. See the full list of supported advanced search query parameters."
//	@Param		hl				query		string								false	"Parameter defines the language to use for the Google search. It's a two-letter language code. (e.g., en for English, es for Spanish, or fr for French). Head to the Google languages page for a full list of supported Google languages."
//	@Param		gl				query		string								false	"Parameter defines the country to use for the Google search. It's a two-letter country code. (e.g., us for the United States, uk for United Kingdom, or fr for France). Head to the Google countries page for a full list of supported Google countries."
//	@Param		date			query		string								false	"Parameter date. last hour=h last day=d last week=w last month=m last year=y"
//	@Param		start			query		int									false	"Parameter defines the result offset. It skips the given number of results. It's used for pagination. (e.g., 0 (default) is the first page of results, 10 is the 2nd page of results, 20 is the 3rd page of results, etc.).Google Local Results only accepts multiples of 20(e.g. 20 for the second page results, 40 for the third page results, etc.) as the start value."
//	@Param		num				query		int									false	"Parameter defines the maximum number of results to return. (e.g., 10 (default) returns 10 results, 40 returns 40 results, and 100 returns 100 results).The use of num may introduce latency, and/or prevent the inclusion of specialized result types. It is better to omit this parameter unless it is strictly necessary to increase the number of results per page."
//	@Param		device			query		string								false	"Parameter defines the device to use to get the results. Currently only supports desktop version."
//	@Param		google_domain	query		string								false	"Parameter defines the Google domain to use. It defaults to google.com. Head to the Google domains page for a full list of supported Google domains."
//	@Param		html			query		string								false	"Parameter html 1 or 0"
//	@Success	200				{object}	base.Resp{data=dto.SearchResp}	"{"code": 200, "data": [...]}"
//	@Router		/api/search [get]
func (e *SearchApi) SearchGet(c *gin.Context) {
	st := time.Now()

	var req dto.SearchReq

	if err := c.ShouldBindQuery(&req); err != nil {
		e.Error(c, err)
		return
	}

	var res dto.SearchResp

	code := service.SerSearchService.SearchHandler(&req, &res)

	df := time.Since(st)

	fmt.Println(df)
	if code != 200 {
		e.Fail(c, code, fmt.Sprintf("code:%d", code))
		return
	}
	e.Ok(c, res)
}

// SearchPost google search
//
//	@Summary	google search
//	@Tags		open
//	@Accept		application/json
//	@Product	application/json
//	@Param		data	body		dto.SearchReq					true	"body"
//	@Success	200		{object}	base.Resp{data=dto.SearchResp}	"{"code": 200, "data": [...]}"
//	@Router		/api/search [post]
func (e *SearchApi) SearchPost(c *gin.Context) {
	st := time.Now()

	var req dto.SearchReq

	if err := c.ShouldBindJSON(&req); err != nil {
		e.Error(c, err)
		return
	}

	var res dto.SearchResp

	code := service.SerSearchService.SearchHandler(&req, &res)

	df := time.Since(st)

	fmt.Println(df)

	if code != 200 {
		e.Fail(c, code, fmt.Sprintf("code:%d", code))
		return
	}
	e.Ok(c, res)
}

var (
	keywords = make([]string, 0)
)

func init() {
	f, err := os.OpenFile("keyword.txt", os.O_RDONLY, 0666)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		keywords = append(keywords, scanner.Text())
	}
}

//	@Summary	Test
//	@Tags		open
//	@Product	application/json
//	@Param		qps				query		int			false	"并发."
//	@Param		loops			query		int			false	"循环次数 1 or 0"
//	@Success	200				{object}	base.Resp{data=map[int]int}	"{"code": 200, "data": [...]}"
//
// @Router		/api/test [get]
func (e *SearchApi) Test(c *gin.Context) {
	strQps := c.Query("qps")
	strLoops := c.Query("loops")
	qps, _ := strconv.Atoi(strQps)
	loops, _ := strconv.Atoi(strLoops)
	errm := make(map[int]int, 0)
	lock := sync.Mutex{}
	total := qps * loops
	wg := &sync.WaitGroup{}
	wg.Add(total)

	for i := 0; i < loops; i++ {
		for j := 0; j < qps; j++ {
			params := &dto.SearchReq{
				Q: keywords[rand.IntN(len(keywords))],
			}
			go func() {
				defer wg.Done()
				var res dto.SearchResp
				code := service.SerSearchService.SearchHandler(params, &res)
				lock.Lock()
				if cnt, ok := errm[code]; ok {
					errm[code] = cnt + 1
				} else {
					errm[code] = 1
				}

				lock.Unlock()
			}()
		}
		time.Sleep(1 * time.Second)
	}
	wg.Wait()
	fmt.Printf("err map: %+v\n", errm)
	e.Ok(c, errm)
}
