package job

import (
	"context"
	"dilu/common/utils"
	"dilu/modules/search/google/handler"
	"net/url"
	"time"

	"github.com/baowk/dilu-core/core"
)

func Init() {
	go IpCheck()
}

func IpCheck() {
	tk := time.NewTicker(time.Second * 5)
	defer tk.Stop()

	for {
		<-tk.C
		tk.Reset(time.Second * 60)
		rcli, err := core.CacheRedis()
		if err != nil {
			continue
		}
		for {
			ip, err := rcli.LPop(context.Background(), handler.Proxy_fail_key).Result()
			if err != nil {
				break
			}
			m := utils.ProxyCheckV2([]string{ip})

			u, err := url.Parse(ip)
			if err != nil {
				continue
			}

			if m[u.Hostname()] {
				handler.BackProxy(ip)
			} else {
				handler.SetFailProxy(ip)
			}
		}
	}
}
