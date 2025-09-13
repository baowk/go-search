package service

import (
	"dilu/modules/search/google/handler"
	"dilu/modules/search/service/dto"

	"github.com/baowk/dilu-core/core/base"
)

type SearchService struct {
	base.BaseService
}

var SerSearchService = SearchService{}

func (e *SearchService) SearchHandler(req *dto.SearchReq, res *dto.SearchResp) int {
	code, err := handler.ToSearch(req, res)
	if err != nil {
		return code
	}
	// 请求数据库
	return 200
}
