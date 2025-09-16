package ua

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/mssola/useragent"
)

func TestUa(t *testing.T) {
	uads := make([]UserAgentData, 0, 100)
	f, err := os.OpenFile("user-agents.json", os.O_RDONLY, 664)
	if err != nil {
		t.Fatal(err)
	}

	err = json.NewDecoder(f).Decode(&uads)
	if err != nil {
		t.Fatal(err)
	}

	var wf *os.File
	wf, err = os.OpenFile("uas.json", os.O_CREATE|os.O_WRONLY, 664)
	if err != nil {
		t.Fatal(err)
	}
	j := json.NewEncoder(wf)

	nuas := make(map[string]*UserAgent, 0)
	for _, uad := range uads {
		//t.Log(uad.UserAgent)
		ua := useragent.New(uad.UserAgent)
		bname, bver := ua.Browser()

		arr := strings.Split(bver, ".")
		var mv string
		if len(arr) > 0 {
			mv = arr[0]
		}

		nua := UserAgent{
			Bot:      ua.Bot(),
			Browser:  bname,
			Mob:      ua.Mobile(),
			Os:       ua.OS(),
			Platform: ua.Platform(),
			UA:       ua.UA(),
			Ver:      bver,
			MainVer:  mv,
		}
		nuas[uad.UserAgent] = &nua
	}
	list := make([]*UserAgent, 0, 1000)
	for _, ua := range nuas {
		list = append(list, ua)
	}
	j.Encode(list)
}
