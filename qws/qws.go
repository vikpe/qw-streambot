package qws

import (
	"errors"
	"fmt"
	"sort"

	"github.com/go-resty/resty/v2"
	"github.com/vikpe/serverstat/qserver/mvdsv"
)

type ServerList []mvdsv.Mvdsv

func GetMvdsvServers() ServerList {
	serversUrl := "https://metaqtv.quake.se/v2/servers/mvdsv"
	resp, err := resty.New().R().SetResult(ServerList{}).Get(serversUrl)

	if err != nil {
		fmt.Println("server fetch error", err.Error())
		return make(ServerList, 0)
	}

	servers := resp.Result().(*ServerList)
	return *servers
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

	return IsSpeccable(server)
}

func IsSpeccable(server mvdsv.Mvdsv) bool {
	if len(server.QtvStream.Url) > 0 {
		return true
	}

	return server.SpectatorSlots.Free > 0 && !RequiresPassword(server.Settings.GetInt("needpass", 0))
}

func RequiresPassword(needpass int) bool {
	if 0 == needpass {
		return false
	}
	const spectatorPasswordBit = 2
	return (needpass & spectatorPasswordBit) > 0
}
