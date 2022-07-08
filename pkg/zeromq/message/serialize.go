package message

import (
	"fmt"

	"github.com/goccy/go-json"
)

type SerializedValue []byte

func Serialize(data any) SerializedValue {
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

func (d SerializedValue) ToString() string {
	var target string
	d.To(&target)
	return target
}

func (d SerializedValue) ToInt() int {
	var target int
	d.To(&target)
	return target
}

func (d SerializedValue) To(target interface{}) {
	Unserialize(d, target)
}
