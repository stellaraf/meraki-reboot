package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"stellar.af/meraki-reboot/util"
)

const WARN string = "https://emojipedia-us.s3.dualstack.us-west-1.amazonaws.com/thumbs/240/apple/237/warning-sign_26a0.png"
const OK string = "https://emojipedia-us.s3.dualstack.us-west-1.amazonaws.com/thumbs/240/apple/237/white-heavy-check-mark_2705.png"

type SlackArgs struct {
	Org     string
	Msg     string
	Title   *string
	Success *bool
}

func SendWebhook(args SlackArgs) error {
	url := util.GetEnv("MERAKI_REBOOT_SLACK_URL")
	if url != "" {
		imageURL := WARN
		success := false
		if args.Success != nil {
			if *args.Success == true {
				success = true
			}
		}

		title := fmt.Sprintf("%s Meraki Reboot: Error", args.Org)
		if success {
			title = fmt.Sprintf("%s Meraki Reboot: Success", args.Org)
			imageURL = OK
		}
		if args.Title != nil {
			title = *args.Title
		}
		// Title
		headerBlock := slack.NewTextBlockObject("mrkdwn", title, false, false)
		headerSection := slack.NewSectionBlock(headerBlock, nil, nil)

		// Body
		msgBlock := slack.NewTextBlockObject("mrkdwn", args.Msg, false, false)
		fields := make([]*slack.TextBlockObject, 0)
		fields = append(fields, msgBlock)

		// Status image
		image := slack.NewImageBlockElement(imageURL, "Status")
		acc := slack.NewAccessory(image)
		sectionBlock := slack.NewSectionBlock(nil, fields, acc)

		// All blocks
		blockset := make([]slack.Block, 0)
		blockset = append(blockset, headerSection, sectionBlock)
		blocks := slack.Blocks{BlockSet: blockset}

		// Webhook data
		msg := slack.WebhookMessage{Text: title, Blocks: &blocks}

		err := slack.PostWebhook(url, &msg)
		if err != nil {
			log.Error("Error sending Slack webhook:\n%s", err)
			return err
		} else {
			log.Debug("Successfully sent Slack webhook")
		}
	} else {
		log.Debug("'SLACK_WEBHOOK_URL' environment variable is not set, skipping...")
	}
	return nil
}
