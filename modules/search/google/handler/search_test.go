package handler

import (
	"dilu/modules/search/service/dto"
	"encoding/json"
	"testing"
)

func TestSearch(t *testing.T) {
	var sc SimpleCookie
	cStr := `{"c": "AEC=AaJma5teb3VjcMt4jeMfTydkv3Zw0yJsWCjsL6Ztu58oFAnnVlFzdyH04Q; NID=525=Q-vL4vpmNcyQY4cAqUJ9u1izrmmb3NEoOUuA6gUNA0wqJHtiE23zx_9aRTMi93HRZg-Unru8YMgNXjUclq4t-KkGsIOgV8c38PYGONOvg1SWjUJXFqb3jHrovARG_fc4bJDUC_7G57VNpIBYQs-VWXzn6STRExXcOG0Nf7XxC9RAcEAHtgChYqzDAXaz-HXYedKdkMAwvcA3NgaAY2P-aLimmTG5q-yRKRvQkdX9mB7JobROqh7Q", "n": 0}`
	err := json.Unmarshal([]byte(cStr), &sc)
	if err != nil {
		t.Error(err)
		return
	}

	_, _, err = Search(
		&dto.SearchReq{
			Q: "hello",
		},
		"http://qgjn4y8nyqvd:rvzcxtuwsdwj@207.228.208.8:5206",
		&sc,
		"android",
	)
	if err != nil {
		t.Error(err)
		return
	}
}
