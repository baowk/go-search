package main

import (
	"bufio"
	"net/url"
	"os"
	"strings"
)

func main() {
	used := make(map[string]struct{})

	f1, err := os.Open("pb1.txt")
	if err != nil {
		return
	}
	defer f1.Close()
	scanner1 := bufio.NewScanner(f1)

	for scanner1.Scan() {
		line := scanner1.Text()
		u, err := url.Parse(line)
		if err != nil {
			continue
		}
		used[u.Hostname()] = struct{}{}
	}

	f2, err := os.Open("pc1.txt")
	if err != nil {
		return
	}
	defer f2.Close()
	scanner2 := bufio.NewScanner(f2)
	for scanner2.Scan() {
		line := scanner2.Text()
		arr := strings.Split(line, ":")
		if len(arr) > 0 {
			used[arr[0]] = struct{}{}
		}

	}

	f, err := os.Open("proxies.txt")
	if err != nil {
		return
	}
	defer f.Close()
	proxies := make(map[string]struct{})
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		u, err := url.Parse(line)
		if err != nil {
			continue
		}
		if _, ok := used[u.Hostname()]; !ok {
			proxies[line] = struct{}{}
		}
	}

	// Convert map keys to slice manually
	list := make([]string, 0, len(proxies))
	for k := range proxies {
		list = append(list, k)
	}

	if err := os.WriteFile("ok.txt", []byte(strings.Join(list, " ")), 0644); err != nil {
		return
	}

}
