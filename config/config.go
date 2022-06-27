package config

import (
	"io/ioutil"
	"os"

	"github.com/goccy/go-json"
)

type StreambotConfig struct {
	Mode               string `json:"mode"`
	MapChangeTimestamp uint   `json:"map_change_timestamp"`
	ServerAddress      string `json:"server_address"`
	SpecAddress        string `json:"spec_address"`
}

func NewFromFile(path string) (StreambotConfig, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return StreambotConfig{}, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config StreambotConfig
	json.Unmarshal(byteValue, &config)

	return config, nil
}
