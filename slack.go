package main

// handles slack API interface for sending webhooks back with responses

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/parnurzeal/gorequest"
)

// Field - struct
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// BotDMPayload - struct for bot DMs
type BotDMPayload struct {
	Token          string       `json:"token,omitempty"`
	Channel        string       `json:"channel,omitempty"`
	Text           string       `json:"text,omitempty"`
	AsUser         bool         `json:"as_user,omitempty"`
	Attachments    []Attachment `json:"attachments,omitempty"`
	IconEmoji      string       `json:"icon_emoji,omitempty"`
	IconURL        string       `json:"icon_url,omitempty"`
	LinkNames      bool         `json:"link_names,omitempty"`
	Mkrdwn         bool         `json:"mrkdwn,omitempty"`
	Parse          string       `json:"parse,omitempty"`
	ReplyBroadcast bool         `json:"reply_broadcast,omitempty"`
	ThreadTS       string       `json:"thread_ts,omitempty"`
	UnfurlLinks    bool         `json:"unfurl_links,omitempty"`
	UnfurlMedia    bool         `json:"unfurl_media,omitempty"`
	Username       string       `json:"username,omitempty"`
}

// Attachment - struct
type Attachment struct {
	Fallback   string   `json:"fallback,omitempty"`
	Color      string   `json:"color,omitempty"`
	PreText    string   `json:"pretext,omitempty"`
	AuthorName string   `json:"author_name,omitempty"`
	AuthorLink string   `json:"author_link,omitempty"`
	AuthorIcon string   `json:"author_icon,omitempty"`
	Title      string   `json:"title,omitempty"`
	TitleLink  string   `json:"title_link,omitempty"`
	Text       string   `json:"text,omitempty"`
	ImageURL   string   `json:"image_url,omitempty"`
	Fields     []*Field `json:"fields,omitempty"`
	Footer     string   `json:"footer,omitempty"`
	FooterIcon string   `json:"footer_icon,omitempty"`
	Timestamp  int64    `json:"ts,omitempty"`
	MarkdownIn []string `json:"mrkdwn_in,omitempty"`
}

// Payload - struct
type Payload struct {
	Parse       string       `json:"parse,omitempty"`
	Username    string       `json:"username,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Text        string       `json:"text,omitempty"`
	LinkNames   string       `json:"link_names,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	UnfurlLinks bool         `json:"unfurl_links,omitempty"`
	UnfurlMedia bool         `json:"unfurl_media,omitempty"`
}

// ReactionPayload - payload to send an emoji reaction to a message
type ReactionPayload struct {
	Token     string `json:"token"`
	Name      string `json:"name"`
	Channel   string `json:"channel"`
	TimeStamp string `json:"timestamp"`
}

const (
	reactionAddURL string = "https://slack.com/api/reactions.add"
)

// AddField - add fields
func (attachment *Attachment) AddField(field Field) *Attachment {
	attachment.Fields = append(attachment.Fields, &field)
	return attachment
}

func redirectPolicyFunc(req gorequest.Request, via []gorequest.Request) error {
	return fmt.Errorf("incorrect token (redirection)")
}

// Send - send message
func Send(webhookURL string, proxy string, payload Payload) []error {
	request := gorequest.New().Proxy(proxy)
	resp, _, err := request.
		Post(webhookURL).
		RedirectPolicy(redirectPolicyFunc).
		Send(payload).
		End()

	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return []error{fmt.Errorf("error sending msg. Status: %v", resp.Status)}
	}

	return nil
}

// Wrangler - wrangle slack webhook calls
func Wrangler(webhookURL string, message string, myChannel string, attachments Attachment) {

	payload := Payload{
		Text:        message,
		Username:    "tagger",
		Channel:     myChannel,
		IconEmoji:   ":spray-paint:",
		Attachments: []Attachment{attachments},
	}
	err := Send(webhookURL, "", payload)
	if len(err) > 0 {
		fmt.Printf("Slack Messaging Error in Wrangler function in slack.go: %s\n", err)
	}
}

// AddReaction - add an emoji reaction to a message (expects proper ReactionPayload struct)
func AddReaction(tagbot TagBot, payload ReactionPayload) error {

	payload.Token = tagbot.SlackBotToken

	jsonStr, err := json.Marshal(&payload)
	if err != nil {
		Logit("Error attempting to marshal struct to json for slack AddReaction", false, "err")
		return err
	}

	req, err := http.NewRequest("POST", reactionAddURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		Logit("Error in http.NewRequest in `AddReaction` in `slack.go`", false, "err")
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+tagbot.SlackBotToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Logit("Error in client.Do in `CreateList` in `trello.go`", false, "err")
		return err
	}
	defer resp.Body.Close()
	return err
}
