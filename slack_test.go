package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/nlopes/slack"
)

const (
	SLACK_HOST_NAME   = "<http://google.com|google.com>"
	EXPECTED_HOSTNAME = "google.com"
)

func TestParseName(t *testing.T) {
	result := parseName(SLACK_HOST_NAME)

	if result != EXPECTED_HOSTNAME {
		t.Error(
			"expected", EXPECTED_HOSTNAME,
			"got", result,
		)
	}

}

func TestFormatMessage(t *testing.T) {
	attr := Attribute{
		CheckCommand: "apt",
		DisplayName:  "apt",
		Name:         "apt",
		State:        3.0,
		CheckTime:    1501416820.43347,
		HostName:     "d0b2c373f2ac",
	}

	expectedAttachment := &slack.Attachment{
		Title:     "Icinga 2 - Dashboard",
		TitleLink: fmt.Sprintf("http://%s/icingaweb2/dashboard", os.Getenv("ICINGA_HOST")),
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
	}

	resultAttachment := formatMessage(attr, SERVICES)

	if !reflect.DeepEqual(expectedAttachment, resultAttachment) {
		t.Error(
			"expected", expectedAttachment,
			"got", resultAttachment,
		)
	}
}
