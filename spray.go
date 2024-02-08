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
func AddSprayCan(e string, paint SprayCans, TagBot TagBot, client *socketmode.Client) (success bool, msg string) {

	var (
		exists bool = false
	)

	// break down "e" into the spray can requested
	tmp := strings.Split(e, " ")
	if len(tmp) < 5 {
		return false, "Invalid command. Use `@tagger add spray can <emoji name (no colons)>`"
	} else {
		e = tmp[4]
	}

	// check if spray can already exists
	for _, sc := range paint {
		if sc.Spray == e {
			exists = true
		}
	}
	if exists {
		return false, "Spray Can " + e + " already exists!"
	}

	// Check slack to see if emoji exists???
	if !ScanEmojiList(e, client) {
		return false, "Emoji " + e + " does not exist on the server!"
	}

	// Add e to spray cans
	paint = append(paint, struct {
		Spray string   `json:"spray"`
		Words []string `json:"words"`
	}{
		Spray: e,
		Words: nil,
	})

	// json marshal and write to file
	err := WriteJSONTagsFile(TagBot.SprayJSONPath, paint)
	if err != nil {
		return false, "Error writing to tags.json: " + err.Error()
	}

	return true, "Spray Can `" + e + "` added!\nNow add some key words!"
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
func ListSprayCans(ev *slackevents.AppMentionEvent, paint SprayCans, TagBot TagBot, client *socketmode.Client) error {

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
		payload.Channel = ev.User
		attachments.Color = "#00ff00"
		attachments.Text = message
		payload.Attachments = append(payload.Attachments, attachments)

		err := WranglerDM(TagBot, payload)
		if err != nil {
			return err
		}

		message = ""
		payload.Attachments = nil
	}

	return nil
}

// SendHelp - send help to a user @tagger help
func SendHelp(user string, TagBot TagBot, client *socketmode.Client) {

	var (
		payload    BotDMPayload
		attachment Attachment
	)

	payload.Text = ""
	payload.Channel = user

	attachment = Attachment{
		Color: "#36a64f",
		Fields: []*Field{
			{
				Title: "tagger Help",
				Value: "tagger is a Slack bot that tags messages with emojis",
				Short: false,
			},
			{
				Title: "Commands",
				Value: "@tagger help` - Get help\n`@tagger list spray cans` - List all tags\n`@tagger add spray can` - Add a tag\n`@tagger delete spray can` - Delete a tag\n`@tagger reload spray cans` - Reload tags.json",
				Short: false,
			},
			{
				Title: "",
				Value: "`@tagger add word` - Add keyword to a spray can (tag)\n`@tagger delete word` - Delete a spray can (tag)",
				Short: false,
			},
		},
	}
	payload.Attachments = append(payload.Attachments, attachment)
	attachment = Attachment{
		Color: "#935DFF",
		Fields: []*Field{
			{
				Title: "Specifics for Adding Words to Spray Cans",
				Value: "`@tagger add word <Spray Can> <new word>`\ne.g.: `@taggerbot add word smile happyness`\nThe <Spray Can> must exist as a real slack emoji.",
				Short: false,
			},
			{
				Title: "Specifics for Adding new Spray Cans",
				Value: "`@tagger add spray can <emoji name (no colons)>`\ne.g.: `@taggerbot add spray can catwave`\nThe <emoji name> must exist as a real slack emoji.",
				Short: false,
			},
		},
	}
	payload.Attachments = append(payload.Attachments, attachment)

	_ = WranglerDM(TagBot, payload)

}

// ScanEmojiList - Find an emoji in the server master emoji list
func ScanEmojiList(emoji string, client *socketmode.Client) bool {

	ServerEmojiList, _ := client.GetEmoji()

	for key := range ServerEmojiList {
		if key == emoji {
			return true
		}
	}

	return false
}
