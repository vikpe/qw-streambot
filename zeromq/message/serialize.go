package message

import "github.com/goccy/go-json"

func Serialize(data any) string {
	dataAsJson, _ := json.Marshal(data)
	return string(dataAsJson)
}

func Unserialize(data string, target interface{}) {
	json.Unmarshal([]byte(data), &target)
}
