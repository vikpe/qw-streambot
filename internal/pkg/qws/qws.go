package qws

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/vikpe/go-qwhub"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
	"github.com/vikpe/streambot/internal/pkg/mtag"
)

func FindPlayer(pattern string) (mvdsv.Mvdsv, error) {
	const minFindLength = 2

	if len(pattern) < minFindLength {
		return mvdsv.Mvdsv{}, errors.New(fmt.Sprintf(`provide at least %d characters.`, minFindLength))
	}

	if !strings.Contains(pattern, "*") {
		pattern = fmt.Sprintf("*%s*", pattern)
	}

	servers := qwhub.NewClient().MvdsvServers(map[string]string{"has_player": pattern})

	if 0 == len(servers) {
		return mvdsv.Mvdsv{}, errors.New(fmt.Sprintf(`player "%s" not found.`, pattern))
	}

	return servers[0], nil
}

func ServerScoreBonus(server mvdsv.Mvdsv) int {
	if !server.Mode.IsXonX() {
		return 0
	} else if !mtag.IsOfficial(server.Settings.Get("matchtag", "")) {
		return 0
	} else if server.Mode.Is1on1() && server.PlayerSlots.Free > 0 {
		return 0
	} else if server.Mode.Is2on2() && server.PlayerSlots.Free > 1 {
		return 0
	} else if server.Mode.Is4on4() && server.PlayerSlots.Free > 2 {
		return 0
	} else {
		return 30
	}
}

func GetBestServer() (mvdsv.Mvdsv, error) {
	servers := qwhub.NewClient().MvdsvServers()

	// add custom score
	for _, server := range servers {
		server.Score += ServerScoreBonus(server)

		if server.Geo.Region != "Europe" {
			server.Score = int(math.Floor(float64(server.Score / 2)))
		}
	}

	sort.Slice(servers, func(i, j int) bool {
		if servers[i].Score == servers[j].Score {
			return servers[i].QtvStream.ID > 0 && servers[j].QtvStream.ID == 0
		}

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
	if server.Geo.Region != "Europe" && len(server.QtvStream.Address) == 0 {
		return false
	} else if server.Mode.IsFortress() {
		return false
	}

	return analyze.IsSpeccable(server)
}
