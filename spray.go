package main

import (
	"strings"

	"github.com/slack-go/slack/slackevents"
)

func CheckTags(ev *slackevents.MessageEvent, tagbot TagBot, paint SprayCans) {

	has, tags := cancontains(paint, strings.ToLower(ev.Text))
	if tagbot.Debug {
		// This generates a lot of logs!!!
		Logit("Evaluating: "+ev.Text, false, "info")
	}

	if has {
		var payload ReactionPayload

		payload.Channel = ev.Channel
		payload.Token = tagbot.SlackBotToken
		payload.TimeStamp = ev.TimeStamp

		for _, tag := range tags {
			payload.Name = tag
			err := AddReaction(tagbot, payload)
			if err != nil {
				Logit("Tagit error catch: "+err.Error(), false, "err")
			}
		}
	}
}

func cancontains(paint SprayCans, e string) (has bool, sprayArray []string) {

	for _, p := range paint {
		for _, t := range p.Words {
			if strings.Contains(e, t) {
				sprayArray = append(sprayArray, p.Spray)
			}
		}
	}

	if len(sprayArray) > 0 {
		return true, sprayArray
	}

	return false, nil
}
