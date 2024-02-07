package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// TagBot - Bot Configuration Options
type TagBot struct {
	SlackHook      string `json:"slackhook"`
	SlackAppToken  string `json:"slackapptoken"`
	SlackBotToken  string `json:"slackbottoken"`
	BotID          string `json:"botid"`
	BotName        string `json:"botname"`
	TeamID         string `json:"teamid"`
	TeamName       string `json:"teamname"`
	LogChannel     string `json:"logchannel"`
	Version        string `json:"version"`
	SprayJSONPath  string `json:"sprayjsonpath"`
	Debug          bool   `json:"debug"`
	TotalSprayCans int
	TotalWords     int
}

// SprayCans - struct for storing search data and emoji tags read in from tag.json
type SprayCans []struct {
	Spray string   `json:"spray"`
	Words []string `json:"words"`
}

// LoadBotConfig - Load Main Bot Configuration JSON
func LoadBotConfig(configPath string) (tagbot TagBot, err error) {
	var fileName string

	if configPath == "" {
		fileName = "config.json"
	}

	file, err := os.Open(fileName)
	if err != nil {
		Logit("error opening config.json file: "+err.Error(), true, "err")
		return tagbot, err
	}

	decoded := json.NewDecoder(file)
	err = decoded.Decode(&tagbot)
	if err != nil {
		Logit("error reading invalid config.json file: "+err.Error(), true, "err")
		return tagbot, err
	}

	if tagbot.Debug {
		fmt.Printf("%+v", tagbot)
	}

	return tagbot, nil
}

// LoadSprayCans - Load tag.json tagger data file
func LoadSprayCans(pathname string) (spray SprayCans, err error) {

	var fileName string

	if pathname == "" {
		fileName = "tags.json"
	} else {
		fileName = pathname
	}

	file, err := os.Open(fileName)
	if err != nil {
		Logit("Error opening "+fileName+"  "+err.Error(), true, "err")
		return spray, err
	}

	decoded := json.NewDecoder(file)
	err = decoded.Decode(&spray)
	if err != nil {
		Logit("Error reading invalid tags.json file: "+err.Error(), true, "err")
		return spray, err
	}

	return spray, nil
}

// WriteTagsFile - Write tag.json tagger data file
func WriteTagsFile(pathname string, spray SprayCans) error {
	var fileName string

	if pathname == "" {
		fileName = "tags.json"
	} else {
		fileName = pathname
	}

	file, err := os.Create(fileName)
	if err != nil {
		Logit("Error creating "+fileName+"  "+err.Error(), true, "err")
		return err
	}

	encoded := json.NewEncoder(file)
	err = encoded.Encode(spray)
	if err != nil {
		Logit("Error writing "+fileName+"  "+err.Error(), true, "err")
		return err
	}

	return nil
}
