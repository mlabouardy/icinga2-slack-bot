package main

import (
	"github.com/mlabouardy/icinga2-slack-bot/config"
	"github.com/nlopes/slack"
)

var (
	icinga2  Icinga2
	config   Config
	slackBot SlackBot
)

func init() {
	config = Config{}
	config.Read()

	icinga2 = Icinga2{
		Host:     config.Icinga2.Host,
		Username: config.Icinga2.Username,
		Password: config.Icinga2.Password,
	}

	slackBot = SlackBot{
		botCommandChannel: make(chan *BotCentral),
		botReplyChannel:   make(chan AttachmentChannel),
		api:               slack.New(config.Slack.Token),
	}
}

func main() {
	slackBot.Connect()
}
