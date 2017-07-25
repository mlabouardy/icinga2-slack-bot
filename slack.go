package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/nlopes/slack"
)

type SlackBot struct {
	api               *slack.Client
	botReplyChannel   chan AttachmentChannel
	botCommandChannel chan *BotCentral
}

type BotCentral struct {
	Channel *slack.Channel
	Event   *slack.MessageEvent
	UserId  string
}

type AttachmentChannel struct {
	Channel      *slack.Channel
	Attachement  *slack.Attachment
	DisplayTitle string
}

type Status struct {
	Label string
	Color string
}

const (
	GREEN   = "#4CAE50"
	RED     = "#F44336"
	ORGANGE = "#FE9700"
	GRAY    = "#607D8B"
	BLUE    = "#03A8F3"
)

var (
	botId          string
	LIST_OF_STATUS = map[float32]Status{
		0: Status{
			Label: "OK",
			Color: GREEN,
		},
		1: Status{
			Label: "WARNING",
			Color: ORGANGE,
		},
		2: Status{
			Label: "CRITICAL",
			Color: RED,
		},
		3: Status{
			Label: "UNKNOWN",
			Color: GRAY,
		},
	}
)

func parseName(name string) string {
	parts := strings.Split(name, "|")
	if len(parts) == 2 {
		return strings.TrimRight(parts[1], ">")
	} else {
		return parts[0]
	}
}

func formatMessage(attr Attribute, objectType ObjectType) *slack.Attachment {
	attachment := &slack.Attachment{
		Title:     "Icinga 2 - Dashboard",
		TitleLink: fmt.Sprintf("http://%s/icingaweb2/dashboard", config.Icinga2.Host),
		Color:     LIST_OF_STATUS[attr.State].Color,
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Command",
				Value: attr.CheckCommand,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Display Name",
				Value: attr.DisplayName,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Name",
				Value: attr.Name,
				Short: true,
			},
			slack.AttachmentField{
				Title: "State",
				Value: LIST_OF_STATUS[attr.State].Label,
				Short: true,
			},
		},
		Footer:     "mlabouardy",
		FooterIcon: "https://yt3.ggpht.com/-VxW-2wCxzHs/AAAAAAAAAAI/AAAAAAAAAAA/fjyskzeA-VA/s900-c-k-no-mo-rj-c0xffffff/photo.jpg",
		Ts:         attr.CheckTime,
	}

	if objectType == SERVICES {
		attachment.Fields = append(attachment.Fields, slack.AttachmentField{
			Title: "Hostname",
			Value: attr.HostName,
			Short: true,
		})
	}

	return attachment
}

func (s *SlackBot) handleBotCommands() {
	commands := map[string]string{
		"help":           "See all the available commands.",
		"service <name>": "Get service <name> status.",
		"host <name>":    "Get host <name> status.",
		"services":       "List all services.",
		"hosts":          "List all hosts.",
	}

	var attachmentChannel AttachmentChannel

	for {
		botChannel := <-s.botCommandChannel
		attachmentChannel.Channel = botChannel.Channel
		commandArray := strings.Fields(botChannel.Event.Text)
		if len(commandArray) >= 2 {
			switch commandArray[1] {
			case "help":
				fields := make([]slack.AttachmentField, 0)
				for k, v := range commands {
					fields = append(fields, slack.AttachmentField{
						Title: "<bot> " + k,
						Value: v,
					})
				}

				attachment := &slack.Attachment{
					Pretext: "The following commands are available:",
					Color:   BLUE,
					Fields:  fields,
				}

				attachmentChannel.Attachement = attachment
				s.botReplyChannel <- attachmentChannel
			case "service":
				if len(commandArray) == 2 {
					name := parseName(commandArray[2])
					res, err := icinga2.ServiceStatus(name)
					if err != nil {
						attachment := &slack.Attachment{
							Text:  err.Error(),
							Color: BLUE,
						}
						attachmentChannel.Attachement = attachment
						s.botReplyChannel <- attachmentChannel
					} else {
						for _, o := range res.Results {
							attr := o.Attrs
							attachmentChannel.Attachement = formatMessage(attr, SERVICES)
							s.botReplyChannel <- attachmentChannel
						}
					}
				} else {
					attachment := &slack.Attachment{
						Text:  "You mean '<service> host <name>'",
						Color: BLUE,
					}
					attachmentChannel.Attachement = attachment
					s.botReplyChannel <- attachmentChannel
				}

			case "host":
				if len(commandArray) == 2 {
					name := parseName(commandArray[2])
					res, err := icinga2.HostStatus(name)
					if err != nil {
						attachment := &slack.Attachment{
							Text:  err.Error(),
							Color: BLUE,
						}
						attachmentChannel.Attachement = attachment
						s.botReplyChannel <- attachmentChannel
					} else {
						for _, o := range res.Results {
							attr := o.Attrs
							attachmentChannel.Attachement = formatMessage(attr, HOSTS)
							s.botReplyChannel <- attachmentChannel
						}
					}
				} else {
					attachment := &slack.Attachment{
						Text:  "You mean '<bot> host <name>'",
						Color: BLUE,
					}
					attachmentChannel.Attachement = attachment
					s.botReplyChannel <- attachmentChannel
				}
			case "services":
				res, err := icinga2.ListServices()
				if err != nil {
					log.Fatal(err)
				} else {
					for _, o := range res.Results {
						attr := o.Attrs
						attachmentChannel.Attachement = formatMessage(attr, SERVICES)
						s.botReplyChannel <- attachmentChannel
					}
				}
			case "hosts":
				res, err := icinga2.ListHosts()
				if err != nil {
					log.Fatal(err)
				} else {
					for _, o := range res.Results {
						attr := o.Attrs
						attachmentChannel.Attachement = formatMessage(attr, HOSTS)
						s.botReplyChannel <- attachmentChannel
					}
				}
			default:
				attachment := &slack.Attachment{
					Text:  "Sorry, I didn't get that!, please type '<bot> help' to get the list of all the available commands",
					Color: BLUE,
				}

				attachmentChannel.Attachement = attachment
				s.botReplyChannel <- attachmentChannel
			}
		} else {
			attachment := &slack.Attachment{
				Text:  "Sorry, I didn't get that!, please type '<bot> help' to get the list of all the available commands",
				Color: BLUE,
			}
			attachmentChannel.Attachement = attachment
			s.botReplyChannel <- attachmentChannel
		}

	}
}

func (s *SlackBot) handleBotReply() {
	for {
		ac := <-s.botReplyChannel
		params := slack.PostMessageParameters{}
		params.AsUser = true
		params.Attachments = []slack.Attachment{*ac.Attachement}
		if _, _, err := s.api.PostMessage(ac.Channel.Name, ac.DisplayTitle, params); err != nil {
			log.Fatal(err)
		}
	}
}

func (s *SlackBot) Connect() {
	rtm := s.api.NewRTM()

	s.api.SetDebug(false)

	go rtm.ManageConnection()
	go s.handleBotCommands()
	go s.handleBotReply()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				botId = ev.Info.User.ID
			case *slack.MessageEvent:
				channelInfo, err := s.api.GetChannelInfo(ev.Channel)
				if err != nil {
					log.Fatal(err)
				}

				botCentral := &BotCentral{
					Channel: channelInfo,
					Event:   ev,
					UserId:  ev.User,
				}

				if ev.Type == "message" && strings.HasPrefix(ev.Text, "<@"+botId+">") {
					s.botCommandChannel <- botCentral
				}
			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())
			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop
			default:
			}
		}
	}
}
