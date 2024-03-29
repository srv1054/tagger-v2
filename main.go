package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/slack-go/slack/socketmode"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func main() {

	var (
		version = "2.0.1"

		attachment Attachment
	)

	// deal with CLI
	v := flag.Bool("v", false, "Show current version")
	cp := flag.String("c", "", "Path to configuration file")
	jp := flag.String("j", "", "Path to SprayCan JSON file")
	flag.Parse()
	if *v {
		fmt.Println("Version: " + version)
		fmt.Println("taggerbot is a slack bot that tags messages with emojis")
		fmt.Println("github.com/srv1054/tagger")
		os.Exit(0)
	}

	configPath := *cp
	sprayPath := *jp

	// Load JSON Configurations
	TagBot, err := LoadBotConfig(configPath)
	if err != nil {
		os.Exit(1)
	}
	if sprayPath == "" {
		if TagBot.SprayJSONPath == "" {
			sprayPath = "tags.json"
		} else {
			sprayPath = TagBot.SprayJSONPath
		}
	}
	Spray, err := LoadSprayCans(sprayPath)
	if err != nil {
		os.Exit(1)
	}
	TagBot.TotalSprayCans = len(Spray)
	for _, v := range Spray {
		TagBot.TotalWords += len(v.Words)
	}

	TagBot.Version = version

	// Check for required variables
	if TagBot.SlackAppToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN must be set in config.json.\n")
		os.Exit(1)
	}
	if !strings.HasPrefix(TagBot.SlackAppToken, "xapp-") {
		fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN must have the prefix \"xapp-\".")
	}
	if TagBot.SlackBotToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must be set in config.json.\n")
		os.Exit(1)
	}
	if !strings.HasPrefix(TagBot.SlackBotToken, "xoxb-") {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}

	// Start the bot
	// Say Hello to slack logging if enabled
	if TagBot.LogChannel != "" && TagBot.SlackHook != "" {
		Wrangler(TagBot.SlackHook, "Tagger `v"+version+"` is starting up", TagBot.LogChannel, attachment)
		Wrangler(TagBot.SlackHook, strconv.Itoa(TagBot.TotalSprayCans)+" Spray Cans loaded via tags.json", TagBot.LogChannel, attachment)
		Wrangler(TagBot.SlackHook, ">"+strconv.Itoa(TagBot.TotalWords)+" Words loaded via tags.json", TagBot.LogChannel, attachment)
	}

	// Setup Slack API
	api := slack.New(
		TagBot.SlackBotToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(TagBot.SlackAppToken),
	)

	// Start Socket Mode
	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	go func() {
		for evt := range client.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				fmt.Println("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				fmt.Println("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				fmt.Println("Connected to Slack with Socket Mode.")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					if TagBot.Debug {
						fmt.Printf("Ignored %+v\n", evt)
					}
					continue
				}

				if TagBot.Debug {
					fmt.Printf("Event received: %+v\n", eventsAPIEvent)
				}

				client.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {

					// Check for mentions of the bot
					case *slackevents.AppMentionEvent:
						if strings.Contains(ev.Text, strings.ToLower("reload tags")) {
							Spray, err = LoadSprayCans(TagBot.SprayJSONPath)
							Logit("Reloading tags.json", false, "info")
							_, _, err := client.PostMessage(ev.Channel, slack.MsgOptionText("Sure, I have reloaded your tags.json file.", false))
							if err != nil {
								Logit("failed posting message: "+err.Error(), false, "err")
							}
						}
						if strings.Contains(ev.Text, strings.ToLower("list spray cans")) {
							_, _, err := client.PostMessage(ev.Channel, slack.MsgOptionText("DM'ing you a list of available Spray Cans (Tags)!", false))
							if err != nil {
								Logit("failed posting message: "+err.Error(), false, "err")
							}
							err = ListSprayCans(ev, Spray, TagBot, client)
							if err != nil {
								Logit("Error listing tags: "+err.Error(), false, "err")
							}
						}
						if strings.Contains(ev.Text, strings.ToLower("add spray can")) {
							success, msg := AddSprayCan(ev.Text, Spray, TagBot, client)
							if !success {
								_, _, _ = client.PostMessage(ev.Channel, slack.MsgOptionText("Failed to add spray can.\n"+msg, false))
							} else {
								_, _, _ = client.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
								// Reload tags that were written to JSON file
								Spray, _ = LoadSprayCans(TagBot.SprayJSONPath)
							}
						}
						if strings.Contains(ev.Text, strings.ToLower("add word")) {
							success, msg := AddWord(ev.Text, Spray, TagBot, client)
							if !success {
								_, _, _ = client.PostMessage(ev.Channel, slack.MsgOptionText("Failed to add new word to spray can.\n"+msg, false))
							} else {
								_, _, _ = client.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
								// Reload tags that were written to JSON file
								Spray, _ = LoadSprayCans(TagBot.SprayJSONPath)
							}
						}
						if strings.Contains(ev.Text, strings.ToLower("delete spray can")) {
							if ev.Channel != TagBot.AllowDeleteFrom {
								_, _, _ = client.PostMessage(ev.Channel, slack.MsgOptionText("You are not in a channel that allows deletions!", false))
								break
							}
							success, msg := DeleteSprayCan(ev.Text, Spray, TagBot, client)
							if !success {
								_, _, _ = client.PostMessage(ev.Channel, slack.MsgOptionText("Failed to add new word to spray can.\n"+msg, false))
							} else {
								_, _, _ = client.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
								// Reload tags that were written to JSON file
								Spray, _ = LoadSprayCans(TagBot.SprayJSONPath)
							}
						}
						if strings.Contains(ev.Text, strings.ToLower("delete word ")) {
							if ev.Channel != TagBot.AllowDeleteFrom {
								_, _, _ = client.PostMessage(ev.Channel, slack.MsgOptionText("You are not in a channel that allows deletions!", false))
								break
							}
							success, msg := DeleteWord(ev.Text, Spray, TagBot, client)
							if !success {
								_, _, _ = client.PostMessage(ev.Channel, slack.MsgOptionText("Failed to add new word to spray can.\n"+msg, false))
							} else {
								_, _, _ = client.PostMessage(ev.Channel, slack.MsgOptionText(msg, false))
								// Reload tags that were written to JSON file
								Spray, _ = LoadSprayCans(TagBot.SprayJSONPath)
							}
						}
						if strings.Contains(ev.Text, strings.ToLower("help")) {
							_, _, err := client.PostMessage(ev.Channel, slack.MsgOptionText("DM'ing you some help!", false))
							if err != nil {
								Logit("failed posting message: "+err.Error(), false, "err")
							}
							SendHelp(ev.User, TagBot, client)
						}

					// Check messages for option to tag them with a spray can
					case *slackevents.MessageEvent:
						if TagBot.Debug {
							fmt.Printf("%v", ev)
						}
						CheckTags(ev, TagBot, Spray)
					}
				default:
					client.Debugf("unsupported Events API event received")
				}
			default:
				if TagBot.Debug {
					fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
				}
			}
		}
	}()

	client.Run()
}
