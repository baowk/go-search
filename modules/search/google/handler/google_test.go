package handler

import (
	"dilu/modules/search/service/dto"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	params := &dto.SearchReq{
		ApiKey: "2426d60c0f3b414ca4e1928dc733d42b",
		Engine: "google",
		Q:      "hello",
		Num:    10,
		Start:  0,
	}
	data, err := readFile("../../../htmls/a.html")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	res := dto.SearchResp{}
	_, _err := ParseHtmlAll(data, params, &res)
	if _err != nil {
		fmt.Println(_err)
		os.Exit(1)
	}
	m, _ := json.Marshal(res)
	fmt.Println(string(m))
}

func readFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := stat.Size()

	bytes := make([]byte, size)
	_, err = file.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
