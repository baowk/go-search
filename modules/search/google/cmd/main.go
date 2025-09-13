package main

import (
	"bufio"
	"dilu/modules/search/google/handler"
	"dilu/modules/search/service/dto"
	"fmt"
	"math/rand/v2"
	"os"
	"sync"
	"time"
)

var keywords = []string{}

func main() {
	f, err := os.OpenFile("keyword.txt", os.O_RDONLY, 0666)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		keywords = append(keywords, scanner.Text())
	}

	var success int
	errm := make(map[int]int, 0)
	lock := sync.Mutex{}
	qps := 1
	for i := 0; i < 2; i++ {
		for j := 0; j < qps; j++ {
			params := &dto.SearchReq{
				Q: keywords[rand.IntN(len(keywords))],
			}
			go func() {
				var res dto.SearchResp
				code, err := handler.ToSearch(params, &res)
				lock.Lock()
				if err != nil {
					//fmt.Println(err)
					if cnt, ok := errm[code]; ok {
						errm[code] = cnt + 1
					} else {
						errm[code] = 1
					}

				} else {
					// j, _ := json.Marshal(res)
					// fmt.Println(string(j))
					success++
				}
				lock.Unlock()
			}()
		}
		time.Sleep(1 * time.Second)
	}
	time.Sleep(10 * time.Second)
	fmt.Println("1s success:", success)
	fmt.Printf("err map: %+v\n", errm)
}
