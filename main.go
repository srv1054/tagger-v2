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
		version = "1.0.2"

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
	Spray, err := LoadSprayCans(sprayPath)
	if err != nil {
		os.Exit(1)
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
	// Say Hello to slack logging
	if TagBot.LogChannel != "" && TagBot.SlackHook != "" {
		cans := len(SprayCans{})
		Wrangler(TagBot.SlackHook, "Tagger `v"+version+"` is starting up", TagBot.LogChannel, attachment)
		Wrangler(TagBot.SlackHook, strconv.Itoa(cans)+" Spray Cans loaded via tags.json", TagBot.LogChannel, attachment)
	}

	api := slack.New(
		TagBot.SlackBotToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(TagBot.SlackAppToken),
	)

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
					// Handle direct messages to the Bot Name Mention
					/* case *slackevents.AppMentionEvent:
					_, _, err := client.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
					if err != nil {
						fmt.Printf("failed posting message: %v", err)
					} */
					// check messages for option to tag them
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
