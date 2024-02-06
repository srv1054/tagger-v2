package tagger

import (
	"strings"

	"github.com/nlopes/slack"
)

// SprayCans - struct for storing search data and emoji tags read in from tag.json
type SprayCans []struct {
	Spray string   `json:"spray"`
	Words []string `json:"words"`
}

// TagIt - check for and tag sumthin
func TagIt(myBot MyBot, paint SprayCans, ev *slack.MessageEvent) {

	has, tag := cancontains(paint, strings.ToLower(ev.Msg.Text))

	if has {
		var payload ReactionPayload

		payload.Channel = ev.Channel
		payload.Name = tag
		payload.Token = myBot.SlackToken
		payload.TimeStamp = ev.Timestamp

		err := AddReaction(myBot, payload)
		if err != nil {
			errTrap(myBot, "Tagit error catch: ", err)
		}
	}

}

func cancontains(paint SprayCans, e string) (has bool, emoji string) {
	for _, p := range paint {
		for _, t := range p.Words {
			if strings.Contains(e, t) {
				return true, p.Spray
			}
		}
	}

	return false, ""
}
