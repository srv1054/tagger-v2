package tagger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
)

// MyBot - Bot Configuration Options
type MyBot struct {
	SlackHook  string `json:"slackhook"`
	SlackToken string `json:"slacktoken"`
	BotID      string `json:"botid"`
	BotName    string `json:"botname"`
	TeamID     string `json:"teamid"`
	TeamName   string `json:"teamname"`
	LogChannel string `json:"logchannel"`
	Version    string `json:"version"`
	ConfigPath string `json:"config"`
	JSONPath   string `json:"json"`
	Debug      bool   `json:"debug"`
}

// ConfigFile - Configuration File options
type ConfigFile struct {
	LogChannel string `json:"logchannel"`
	Debug      bool   `json:"debug"`
}

// LoadBotConfig - Load Main Bot Configuration TOML
func LoadBotConfig(myBot MyBot) (tmpBot ConfigFile, err error) {
	var fileName string

	if myBot.ConfigPath == "" {
		fileName = "config.json"
	} else {
		if runtime.GOOS == "windows" {
			fileName = myBot.ConfigPath + "\\config.json"
		} else {
			fileName = myBot.ConfigPath + "/config.json"
		}
	}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening config.json file: " + err.Error() + ".  Could not find path for " + fileName)
		return tmpBot, err
	}

	decoded := json.NewDecoder(file)
	err = decoded.Decode(&tmpBot)
	if err != nil {
		fmt.Println("Error reading invalid config.json file: " + err.Error())
		return tmpBot, err
	}

	if tmpBot.Debug {
		fmt.Printf("%+v", tmpBot)
	}

	return tmpBot, nil
}

// LoadSprayCans - Load tag.json tagger data file
func LoadSprayCans(pathname string) (spray SprayCans, err error) {
	var fileName string

	if pathname == "" {
		fileName = "tags.json"
	} else {
		if runtime.GOOS == "windows" {
			fileName = pathname + "\\tags.json"
		} else {
			fileName = pathname + "/tags.json"
		}
	}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening tags.json file: " + err.Error() + ".   Could not find path for " + fileName)
		return spray, err
	}

	decoded := json.NewDecoder(file)
	err = decoded.Decode(&spray)
	if err != nil {
		fmt.Println("Error reading invalid tags.json file: " + err.Error())
		return spray, err
	}

	return spray, nil
}

// errTrap - Generic error handling function
func errTrap(myBot MyBot, message string, err error) {
	var attachments Attachment

	if myBot.Debug {
		fmt.Println(message + "(" + err.Error() + ")")
	}
	attachments.Color = "#ff0000"
	attachments.Text = err.Error()
	LogToSlack(message, myBot, attachments)
}
