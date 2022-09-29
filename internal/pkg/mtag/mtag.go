package mtag

import (
	"fmt"
	"strings"

	"github.com/vikpe/wildcard"
)

var officialTags = []string{
	"qwdl", "qwduel",
	"getquad", "gq",
	"kombat",
}

func IsOfficial(matchtag string) bool {
	matchtag = strings.TrimSpace(matchtag)

	if "" == matchtag {
		return false
	}

	words := strings.SplitN(strings.ToLower(matchtag), " ", 2)

	for _, tag := range officialTags {
		tagPattern := fmt.Sprintf("*%s*", tag)
		if wildcard.Match(tagPattern, words[0]) {
			return true
		}
	}

	return false
}
