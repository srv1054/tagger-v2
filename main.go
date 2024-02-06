package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack/socketmode"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func main() {

	var (
		version = "0.0.5"
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

	// Check for required variables
	if TagBot.SlackAppToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN must be set in config.json.\n")
		os.Exit(1)
	}
	if !strings.HasPrefix(TagBot.SlackAppToken, "xapp-") {
		fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN must have the prefix \"xapp-\".")
	}
	if TagBot.SlackBotToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must be set.\n")
		os.Exit(1)
	}
	if !strings.HasPrefix(TagBot.SlackBotToken, "xoxb-") {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}

	// Start the bot
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
					fmt.Printf("Ignored %+v\n", evt)

					continue
				}

				fmt.Printf("Event received: %+v\n", eventsAPIEvent)

				client.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					/* case *slackevents.AppMentionEvent:
						_, _, err := client.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
						if err != nil {
							fmt.Printf("failed posting message: %v", err)
						}
					case *slackevents.MemberJoinedChannelEvent:
						fmt.Printf("user %q joined to channel %q", ev.User, ev.Channel) */
					case *slackevents.MessageEvent:
						fmt.Printf("%v", ev)
						CheckTags(ev)
						// message of any kind was received
						// call function to check for matching words
						// if found tag with emoji
					}
				default:
					client.Debugf("unsupported Events API event received")
				}
			/* case socketmode.EventTypeInteractive:
				callback, ok := evt.Data.(slack.InteractionCallback)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)

					continue
				}

				fmt.Printf("Interaction received: %+v\n", callback)

				var payload interface{}

				switch callback.Type {
				case slack.InteractionTypeBlockActions:
					// See https://api.slack.com/apis/connections/socket-implement#button

					client.Debugf("button clicked!")
				case slack.InteractionTypeShortcut:
				case slack.InteractionTypeViewSubmission:
					// See https://api.slack.com/apis/connections/socket-implement#modal
				case slack.InteractionTypeDialogSubmission:
				default:

				}

				client.Ack(*evt.Request, payload)
			case socketmode.EventTypeSlashCommand:
				cmd, ok := evt.Data.(slack.SlashCommand)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)

					continue
				}

				client.Debugf("Slash command received: %+v", cmd)

				payload := map[string]interface{}{
					"blocks": []slack.Block{
						slack.NewSectionBlock(
							&slack.TextBlockObject{
								Type: slack.MarkdownType,
								Text: "foo",
							},
							nil,
							slack.NewAccessory(
								slack.NewButtonBlockElement(
									"",
									"somevalue",
									&slack.TextBlockObject{
										Type: slack.PlainTextType,
										Text: "bar",
									},
								),
							),
						),
					},
				}

				client.Ack(*evt.Request, payload) */
			default:
				fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
			}
		}
	}()

	client.Run()
}
