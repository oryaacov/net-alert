package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//Configuration the netAlert configuration struct
type Configuration struct {
	DB struct {
		ConnectionString string `json:"connectionString"`
		Driver           string `json:"driver"`
	} `json:"DB"`
	Log struct {
		LogFilePath  string `json:"logFilePath"`
		LogLevel     string `json:"logLevel"`
		LogToConsole bool   `json:"logToConsole"`
	} `json:"Log"`
	Sniffer struct {
		Device      string `json:"Device"`
		Promiscuous bool   `json:"Promiscuous"`
		SnapshotLen int32  `json:"SnapshotLen"`
		Timeout     int    `json:"Timeout"`
	} `json:"Sniffer"`
	WebServer struct {
		AllowedHeaders      string `json:"AllowedHeaders"`
		AllowedMethods      string `json:"AllowedMethods"`
		SiteURL             string `json:"SiteUrl"`
		StaticFilesLocation string `json:"StaticFilesLocation"`
		URL                 string `json:"URL"`
		Port                int    `json:"port"`
	} `json:"WebServer"`
}

//ReadConfigutionFromFile reads the netAlert.json configuration file into a struct
func ReadConfigutionFromFile(path string) *Configuration {
	var config Configuration
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("can not find configuration file")
	}
	if bytes, err := ioutil.ReadFile(path); err != nil {
		panic(err)
	} else {
		if jsnError := json.Unmarshal(bytes, &config); jsnError != nil {
			panic(jsnError)
		}
	}
	return &config
}
