package config

import (
	"encoding/json"
	"os"
)

var Conf = configuration{}

type configuration struct {
	LocalUrl string `json:"localurl"`
}

func init() {
	file, _ := os.Open("conf/config.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	err := decoder.Decode(&Conf)

	if err != nil {
		panic("parse config file error")
	}
}
