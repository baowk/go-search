package handler

import (
	"context"

	"github.com/baowk/dilu-core/core"
)

const (
	proxy_ok_key   = "google:proxies:ok"
	proxy_fail_key = "google:proxies:fail"
)

func GetProxy() (string, error) {
	rCli, err := core.CacheRedis()
	if err != nil {
		return "", err
	}

	return rCli.LPop(context.Background(), proxy_ok_key).Result()

}

func BackProxy(proxyUrl string) error {
	if proxyUrl == "" {
		return nil
	}
	rCli, err := core.CacheRedis()
	if err != nil {
		return err
	}
	return rCli.RPush(context.Background(), proxy_ok_key, proxyUrl).Err()
}

func SetFailProxy(proxyUrl string) error {
	if proxyUrl == "" {
		return nil
	}
	rCli, err := core.CacheRedis()
	if err != nil {
		return err
	}
	return rCli.RPush(context.Background(), proxy_fail_key, proxyUrl).Err()
}
