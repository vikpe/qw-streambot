package qws

import (
	"errors"
	"fmt"
	"sort"

	"github.com/go-resty/resty/v2"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
)

type ServerList []mvdsv.Mvdsv

func GetMvdsvServers() ServerList {
	return GetMvdsvServersByQueryParams(nil)
}

func GetMvdsvServersByQueryParams(queryParams map[string]string) ServerList {
	serversUrl := "https://metaqtv.quake.se/v2/servers/mvdsv"
	resp, err := resty.New().R().SetResult(ServerList{}).SetQueryParams(queryParams).Get(serversUrl)

	if err != nil {
		fmt.Println("server fetch error", err.Error())
		return make(ServerList, 0)
	}

	servers := resp.Result().(*ServerList)
	return *servers
}

func FindPlayer(name string) (mvdsv.Mvdsv, error) {
	const minFindLength = 2

	if len(name) < minFindLength {
		return mvdsv.Mvdsv{}, errors.New(fmt.Sprintf(`provide at least %d characters.`, minFindLength))
	}

	queryParams := map[string]string{"has_player": name}
	servers := GetMvdsvServersByQueryParams(queryParams)

	if 0 == len(servers) {
		return mvdsv.Mvdsv{}, errors.New(fmt.Sprintf(`player "%s" not found.`, name))
	}

	return servers[0], nil
}

func GetBestServer() (mvdsv.Mvdsv, error) {
	servers := GetMvdsvServers()

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

	return analyze.IsSpeccable(server)
}
