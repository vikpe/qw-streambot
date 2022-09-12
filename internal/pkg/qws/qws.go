package qws

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
)

func GetMvdsvServers(queryParams map[string]string) []mvdsv.Mvdsv {
	serversUrl := "https://hubapi.quakeworld.nu/v2/servers/mvdsv"
	resp, err := resty.New().R().SetResult([]mvdsv.Mvdsv{}).SetQueryParams(queryParams).Get(serversUrl)

	if err != nil {
		fmt.Println("server fetch error", err.Error())
		return make([]mvdsv.Mvdsv, 0)
	}

	servers := resp.Result().(*[]mvdsv.Mvdsv)
	return *servers
}

func FindPlayer(pattern string) (mvdsv.Mvdsv, error) {
	const minFindLength = 2

	if len(pattern) < minFindLength {
		return mvdsv.Mvdsv{}, errors.New(fmt.Sprintf(`provide at least %d characters.`, minFindLength))
	}

	if !strings.Contains(pattern, "*") {
		pattern = fmt.Sprintf("*%s*", pattern)
	}

	servers := GetMvdsvServers(map[string]string{"has_player": pattern})

	if 0 == len(servers) {
		return mvdsv.Mvdsv{}, errors.New(fmt.Sprintf(`player "%s" not found.`, pattern))
	}

	return servers[0], nil
}

func GetBestServer() (mvdsv.Mvdsv, error) {
	servers := GetMvdsvServers(nil)

	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Score > servers[j].Score
	})

	for _, server := range servers {
		if IsRelevantServer(server) {
			return server, nil
		}
	}

	return mvdsv.Mvdsv{}, errors.New("no server found")
}

func IsRelevantServer(server mvdsv.Mvdsv) bool {
	if server.Geo.Region == "South America" {
		return false
	}

	if server.Mode.IsFortress() {
		return false
	}

	return analyze.IsSpeccable(server)
}
