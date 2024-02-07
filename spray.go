package main

import (
	"strings"

	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// CheckTags - checks for tags in a message and adds reactions
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

// cancontains - checks if a string contains any of the words in a spray can
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

// AddSprayCan - adds a spray can to the JSON
func AddSprayCan(e string, paint SprayCans) error {
	// Add a spray can
	// AddSprayCan(ev.Text, Spray)
	return nil
}

// DeleteSprayCan - deletes a spray can from the JSON
func DeleteSprayCan(e string, paint SprayCans) error {
	// Delete a spray can
	// DeleteSprayCan(ev.Text, Spray)
	return nil
}

// AddWord - adds a word to a spray can
func AddWord(e string, paint SprayCans) error {
	// Add a word to a spray can
	// AddWord(ev.Text, Spray)
	return nil
}

// DeleteWord - deletes a word from a spray can
func DeleteWord(e string, paint SprayCans) error {
	// Delete a word from a spray can
	// DeleteWord(ev.Text, Spray)
	return nil
}

// ListTags - lists all the spray cans and their words
func ListSprayCans(ev *slackevents.AppMentionEvent, paint SprayCans, tagbot TagBot, client *socketmode.Client) error {

	var (
		payload     BotDMPayload
		message     string = ""
		hmessage    string = ""
		attachments Attachment
	)

	payload.Attachments = nil

	for _, p := range paint {
		hmessage = "Keywords for tag :" + p.Spray + ":\n"
		for _, w := range p.Words {
			message = message + w + "\n"
		}

		payload.Text = hmessage
		payload.Channel = ev.Channel
		attachments.Color = "#00ff00"
		attachments.Text = message
		payload.Attachments = append(payload.Attachments, attachments)

		err := WranglerDM(tagbot, payload)
		if err != nil {
			return err
		}

		message = ""
		payload.Attachments = nil
	}

	return nil
}
