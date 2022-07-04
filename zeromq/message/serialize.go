package message

import (
	"fmt"

	"github.com/goccy/go-json"
)

func Serialize(data any) []byte {
	dataAsJson, err := json.Marshal(data)

	if err != nil {
		fmt.Println("Serialize error", data, err)
	}

	return dataAsJson
}

func Unserialize(data []byte, target interface{}) {
	err := json.Unmarshal(data, &target)
	if err != nil {
		fmt.Println("Unserialize error", data, err)
	}
}
