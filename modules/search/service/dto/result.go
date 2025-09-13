package dto

type SearchResp struct {
	//CostTime          float64                    `json:"cost_time"`          //
	SearchMetadata    SearchMetadata             `json:"search_metadata"`    //
	SearchInformation []GetSearchInformationItem `json:"search_information"` //
	SearchParameters  SearchReq                  `json:"search_parameters"`  //
	URL               string                     `json:"url"`                //
	RawHtml           string                     `json:"raw_html"`           // 原始html
	//LocalResults      []LocalResultsItem         `json:"local_results"`      //
	//LocalMap          []LocalMapItem             `json:"local_map"`          //
	OrganicResults []OrganicResultsItem `json:"organic_results"` //
	AlsoAsk        []AlsoAskItem        `json:"also_ask"`        //
	//KnowledgeGraph  []KnowledgeGraphItem `json:"knowledge_graph"`  //
	RelatedSearches RelatedSearchesData `json:"related_searches"` //
	Pagination      GetSearchPagination `json:"pagination"`       //
}

type SearchMetadata struct {
	RawHtmlFile    string  `json:"raw_html_file"`
	XrayHtmlFile   string  `json:"xray_html_file"`
	TotalTimeTaken float64 `json:"total_time_taken"`
	ID             string  `json:"id"`
	JSONEndpoint   string  `json:"json_endpoint"`
	CreatedAt      string  `json:"created_at"`
	ProcessedAt    string  `json:"processed_at"`
	GoogleURL      string  `json:"google_url"`
}

type GetSearchInformationItem struct {
	QueryDisplayed string `json:"query_displayed"`
}

type OrganicResultsItem struct {
	Title            string            `json:"title"`
	Description      string            `json:"description,omitempty"`
	Date             string            `json:"date"`
	Url              string            `json:"url"`
	OriginSite       string            `json:"origin_site"`
	Position         int               `json:"position"`
	OriginNavigation string            `json:"origin_navigation"`
	SourceLogo       string            `json:"source_logo"`
	Img              map[string]string `json:"img,omitempty"`
	Detailed         string            `json:"detailed,omitempty"`
	//Favicon          string            `json:"favicon"`
	//Rank        int               `json:"rank"`
}

type AlsoAskItem struct {
	Question string `json:"question"`
}

type KnowledgeGraphInlineImagesItem struct {
	Img       string `json:"img"`
	TitleLink string `json:"title_link"`
}

type RelatedSearchesDataItem struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}
type RelatedSearchesData struct {
	Data []RelatedSearchesDataItem `json:"data"`
}

type GetSearchPagination struct {
	Current    int               `json:"current"`
	Next       string            `json:"next"`
	OtherPages map[string]string `json:"other_pages"`
}
