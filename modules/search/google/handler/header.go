package handler

import (
	"context"
	"encoding/json"

	"github.com/baowk/dilu-core/core"
)

const (
	cookies_key = "serp:google:cookies"

	cookies_fail_key = "serp:google:cookies:fail"
)

func GetReqHeader() (*SimpleCookie, error) {
	rCli, err := core.CacheRedis()
	if err != nil {
		return nil, err
	}

	strC, err := rCli.LPop(context.Background(), cookies_key).Result()

	if err != nil {
		return nil, err
	}
	var sc SimpleCookie
	err = json.Unmarshal([]byte(strC), &sc)

	return &sc, err
}

func BackReqHeader(header *SimpleCookie) error {
	if header == nil {
		return nil
	}
	// if header.N > 39 {
	// 	return nil
	// }
	rCli, err := core.CacheRedis()
	if err != nil {
		return err
	}

	header.N = header.N + 1
	data, _ := json.Marshal(header)
	return rCli.RPush(context.Background(), cookies_key, string(data)).Err()
}

func SetFailReqHeader(header *SimpleCookie) error {
	if header == nil {
		return nil
	}
	// if header.N > 39 {
	// 	return nil
	// }
	rCli, err := core.CacheRedis()
	if err != nil {
		return err
	}

	header.N = header.N + 1
	data, _ := json.Marshal(header)
	return rCli.RPush(context.Background(), cookies_fail_key, string(data)).Err()
}

type SimpleCookie struct {
	C string `json:"c"`
	N int    `json:"n"`
}
