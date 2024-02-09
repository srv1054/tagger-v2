package main

import (
	"strings"
	"unicode"

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

	// Check slack to see if emoji exists
	// This feature does work, however the Slack API does NOT return emojis that are "built-in" to the server, only custom user emojis
	// This is a limitation of the Slack API, so this feature is commented out for now
	//if !ScanEmojiList(e, client) {
	//return false, "Emoji `" + e + "` does not exist on the server!"
	//}

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

	Logit("Spray Can `"+e+"` added to tags.json", false, "info")
	return true, "Spray Can `" + e + "` added!\nNow add some key words!"
}

// AddWord - adds a word to a spray can
func AddWord(e string, paint SprayCans, TagBot TagBot, client *socketmode.Client) (success bool, msg string) {

	var (
		exists   bool = false
		sprayCan string
		word     string
	)

	// break down "e" into the word requested
	tmp := strings.Split(e, " ")
	if len(tmp) < 5 {
		return false, "Invalid command. Use `@tagger add word <spray can> \"<word>\"`"
	} else {
		sprayCan = tmp[3]
		// remove the quotes from the word
		word = strings.Trim(tmp[4], "\"")
	}

	/* Validation
	 - We need to do lots of validation here on what is trying to be added.   Adding simple words or a space can cause tagger to
		tag everything in the channel all the time, which could be bad or overloading.
	 - Should not be NaughtyWords list (strip spaces before compare) */
	for _, nw := range NaughtyWords {
		if strings.ToLower(strings.ReplaceAll(word, " ", "")) == nw {
			return false, "Common word `" + word + "` cannot be added to a Spray Can `" + sprayCan + "`"
		}
	}
	if ContainsOnlySpaces(word) {
		return false, "Word cannot be all spaces!"
	}
	if len(word) < 4 {
		return false, "Words must be greater then 3 characters, sorry `" + word + "` won't work!"
	}
	if word == "" {
		return false, "Word cannot be blank!"
	}

	// find the specified spray can and then validate the word doens't already exist
	for _, sc := range paint {
		if sc.Spray == sprayCan {
			for _, w := range sc.Words {
				if w == word {
					exists = true
				}
			}
		}
	}
	if exists {
		return false, "The word `" + word + "` already exists in Spray Can `" + sprayCan + "`!"
	}

	// Validate that the Spray Can already exists
	exists = false
	for _, sc := range paint {
		if sc.Spray == sprayCan {
			exists = true
		}
	}
	if !exists {
		return false, "Spray Can `" + sprayCan + "` does not exist!"
	}

	// Add word to spray can
	for i, sc := range paint {
		if sc.Spray == sprayCan {
			paint[i].Words = append(paint[i].Words, word)
		}
	}

	// json marshal and write to file
	err := WriteJSONTagsFile(TagBot.SprayJSONPath, paint)
	if err != nil {
		return false, "Error writing to tags.json: " + err.Error()
	}

	Logit("Word `"+word+"` added to Spray Can `"+sprayCan+"` in tags.json", false, "info")
	return true, "Word `" + word + "` added to Spray Can `" + sprayCan + "`!"
}

// DeleteSprayCan - deletes a spray can from the JSON
func DeleteSprayCan(e string, paint SprayCans) error {
	// Delete a spray can
	// DeleteSprayCan(ev.Text, Spray)
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
				Value: "`@tagger add spray word` - Add keyword to a spray can (tag)\nYou must specify an existing Spray Can\nWord *must* be in \"\" quotation marks to allow for spaces.\n`@tagger delete word` - Delete a spray can (tag)",
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

// ContainsOnlySpaces - checks if a string contains only spaces
func ContainsOnlySpaces(s string) bool {
	for _, char := range s {
		if !unicode.IsSpace(char) {
			return false
		}
	}
	return true
}
