package dto

import (
	"fmt"
	"net/url"
	"strconv"
)

type SearchReq struct {
	Engine       string `json:"engine" form:"engine"`                         // Parameter engine Currently only supports google
	Q            string `json:"q" form:"q" binding:"required"`                // Parameter defines the query you want to search. You can use anything that you would use in a regular Google search. e.g. inurl:, site:, intitle:. We also support advanced search query parameters such as as_dt and as_eq. See the full list of supported advanced search query parameters.
	Hl           string `json:"hl,omitempty" form:"hl"`                       // Parameter defines the language to use for the Google search. It's a two-letter language code. (e.g., en for English, es for Spanish, or fr for French). Head to the Google languages page for a full list of supported Google languages.
	Gl           string `json:"gl,omitempty" form:"gl"`                       // Parameter defines the country to use for the Google search. It's a two-letter country code. (e.g., us for the United States, uk for United Kingdom, or fr for France). Head to the Google countries page for a full list of supported Google countries.
	Location     string `json:"location,omitempty" form:"location"`           // Parameter defines from where you want the search to originate. If several locations match the location requested, we'll pick the most popular one. Head to the /locations.json API if you need more precise control. The location and uule parameters can't be used together. It is recommended to specify location at the city level in order to simulate a real user’s search. If location is omitted, the search may take on the location of the proxy.
	Time         string `json:"time,omitempty" form:"time"`                   // Parameter time last hour=h last day=d last week=w last month=m last year=y
	Start        int    `json:"start" form:"start"`                           // Parameter defines the result offset. It skips the given number of results. It's used for pagination. (e.g., 0 (default) is the first page of results, 10 is the 2nd page of results, 20 is the 3rd page of results, etc.).Google Local Results only accepts multiples of 20(e.g. 20 for the second page results, 40 for the third page results, etc.) as the start value.
	Num          int    `json:"num,omitempty" form:"num"`                     // Parameter defines the maximum number of results to return. (e.g., 10 (default) returns 10 results, 40 returns 40 results, and 100 returns 100 results).The use of num may introduce latency, and/or prevent the inclusion of specialized result types. It is better to omit this parameter unless it is strictly necessary to increase the number of results per page.
	Device       string `json:"device,omitempty" form:"device"`               // Parameter defines the device to use to get the results. Currently only supports desktop version.
	ApiKey       string `json:"api_key,omitempty"`                            // Parameter api_key open search required
	GoogleDomain string `json:"google_domain,omitempty" form:"google_domain"` // Parameter defines the Google domain to use. It defaults to google.com. Head to the Google domains page for a full list of supported Google domains.
	Tbs          string `json:"tbs,omitempty"  form:"tbs"`                    // Parameter defines advanced search parameters that aren't possible in the regular query field. (e.g., advanced search for patents, dates, news, videos, images, apps, or text contents). (to be searched)
	Safe         string `json:"safe,omitempty"  form:"safe"`                  // Parameter defines the level of filtering for adult content. It can be set to active or off, by default Google will blur explicit content. off or active
	Filter       string `json:"filter,omitempty"  form:"filter"`              // Parameter defines the exclusion of results from an auto-corrected query when the original query is spelled wrong. It can be set to 1 to exclude these results, or 0 to include them (default). Note that this parameter may not prevent Google from returning results for an auto-corrected query if no other results are available. 1 or 0
	Nfpr         string `json:"nfpr,omitempty" form:"nfpr"`                   // Parameter defines if the filters for 'Similar Results' and 'Omitted Results' are on or off. It can be set to 1 (default) to enable these filters, or 0 to disable these filters. 1 or 0
	NoCache      bool   `json:"no_cache,omitempty" form:"no_cache"`           // Parameter will force SerpApi to fetch the Google results even if a cached version is already present. A cache is served only if the query and all parameters are exactly the same. Cache expires after 1h. Cached searches are free, and are not counted towards your searches per month. It can be set to false (default) to allow results from the cache, or true to disallow results from the cache. no_cache and async parameters should not be used together. true or false
	FetchMode    string `json:"fetch_mode,omitempty" form:"fetch_mode"`       // Parameter fetch_mode Currently only supports static
	Cr           string `json:"cr,omitempty"  form:"cr"`                      // Parameter defines one or multiple countries to limit the search to. It uses country{two-letter upper-case country code} to specify countries and cr="countryCN|countryJP"
	Lr           string `json:"lr,omitempty" form:"lr"`                       // Parameter defines one or multiple languages to limit the search to. It uses lang_{two-letter language code} to specify languages and lr = "lang_de|lang_ja"
	Ludocid      string `json:"ludocid,omitempty" form:"ludocid"`             // Parameter defines the id (CID) of the Google My Business listing you want to scrape. Also known as Google Place ID.
	Lsig         string `json:"lsig,omitempty" form:"lsig"`                   // Parameter that you might have to use to force the knowledge graph map view to show up. You can find the lsig ID by using our Local Pack API or Google Local API.lsig ID is also available via a redirect Google uses within Google My Business.
	Kgmid        string `json:"kgmid,omitempty" form:"kgmid"`                 // Parameter defines the id (KGMID) of the Google Knowledge Graph listing you want to scrape. Also known as Google Knowledge Graph ID. Searches with kgmid parameter will return results for the originally encrypted search parameters. For some searches, kgmid may override all other parameters except start, and num parameters.
	Ibp          string `json:"ibp,omitempty" form:"ibp"`                     // Parameter is responsible for rendering layouts and expansions for some elements (e.g., gwp;0,7 to expand searches with ludocid for expanded knowledge graph).
	Uds          string `json:"uds,omitempty" form:"uds"`                     // Parameter enables to filter search. It's a string provided by Google as a filter. uds values are provided under the section: filters with uds, q and serpapi_link values provided for each filter.
	Html         string `json:"html,omitempty" form:"html"`                   // Parameter html 1 or 0
	Tbm          string `json:"tbm" form:"tbm"`                               // Parameter defines the type of search you want to do. It can be set to: lcl:Google Local API
	Udm          int    `form:"udm" json:"udm"`
	SkType       string `form:"sk_type" json:"sk_type"` //
}

func (s *SearchReq) GetGoogleUrl() string {
	return fmt.Sprintf("https://www.google.com/search?%s", s.ToString())
}

// func (s *SearchReq) ToShowGoogleUrl() string {
// 	return strings.ReplaceAll(s.GetGoogleUrl(), "\u0026", "&")
// }

func (s *SearchReq) ToString() string {
	params := url.Values{}
	//fmt.Println("SearchReq:", s.Q)
	// Always add the query
	if s.Q != "" {
		params.Set("q", s.Q)
	}

	// Add optional parameters
	if s.Hl != "" {
		params.Set("hl", s.Hl)
	}
	if s.Gl != "" {
		params.Set("gl", s.Gl)
	}
	if s.Location != "" {
		params.Set("location", s.Location)
	}
	if s.Time != "" {
		params.Set("time", s.Time)
	}
	if s.Start != 0 {
		params.Set("start", strconv.Itoa(s.Start))
	}
	if s.Num != 0 {
		params.Set("num", strconv.Itoa(s.Num))
	}
	if s.Device != "" {
		params.Set("device", s.Device)
	}
	// if s.GoogleDomain != "" {
	// 	params.Set("google_domain", s.GoogleDomain)
	// }
	if s.Tbs != "" {
		params.Set("tbs", s.Tbs)
	}
	if s.Safe != "" {
		params.Set("safe", s.Safe)
	}
	if s.Filter != "" {
		params.Set("filter", s.Filter)
	}
	if s.Nfpr != "" {
		params.Set("nfpr", s.Nfpr)
	}
	if s.NoCache {
		params.Set("no_cache", "true")
	}
	if s.FetchMode != "" {
		params.Set("fetch_mode", s.FetchMode)
	}
	if s.Cr != "" {
		params.Set("cr", s.Cr)
	}
	if s.Lr != "" {
		params.Set("lr", s.Lr)
	}
	if s.Ludocid != "" {
		params.Set("ludocid", s.Ludocid)
	}
	if s.Lsig != "" {
		params.Set("lsig", s.Lsig)
	}
	if s.Kgmid != "" {
		params.Set("kgmid", s.Kgmid)
	}
	if s.Ibp != "" {
		params.Set("ibp", s.Ibp)
	}
	if s.Uds != "" {
		params.Set("uds", s.Uds)
	}
	if s.Tbm != "" && s.Udm != 0 {
		params.Set("tbm", s.Tbm)
	} else if s.Tbm != "" {
		params.Set("tbm", s.Tbm)
	} else if s.Udm != 0 {
		params.Set("udm", strconv.Itoa(s.Udm))
	}
	if s.SkType != "" && s.Device != "" { //搜索google应用商店
		params.Set("device", s.Device)
	}
	return params.Encode()
}
