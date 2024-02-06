package main

import (
	"strings"

	"github.com/nlopes/slack"

	"github.com/slack-go/slack/slackevents"
)

func CheckTags(ev *slackevents.MessageEvent) {

}

// TagIt - check for and tag sumthin
func TagIt(tagbot TagBot, paint SprayCans, ev *slack.MessageEvent) {

	has, tag := cancontains(paint, strings.ToLower(ev.Msg.Text))

	if has {
		var payload ReactionPayload

		payload.Channel = ev.Channel
		payload.Name = tag
		payload.Token = tagbot.SlackBotToken
		payload.TimeStamp = ev.Timestamp

		err := AddReaction(tagbot, payload)
		if err != nil {
			Logit("Tagit error catch: "+err.Error(), false, "err")
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
