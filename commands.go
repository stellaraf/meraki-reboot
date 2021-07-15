package main

import (
	"fmt"

	"github.com/mkideal/cli"
	log "github.com/sirupsen/logrus"
	"stellar.af/meraki-reboot/meraki"
	"stellar.af/meraki-reboot/util"
)

type argT struct {
	cli.Helper
	OrgName     string `cli:"o,org" usage:"Meraki Organization Name"`
	NetworkName string `cli:"n,network" usage:"Meraki Network Name"`
	Exclusions  string `cli:"e,exclusions" usage:"Comma-separated list of tags to exclude from the results"`
}

type devicesT struct {
	cli.Helper
}

type slackT struct {
	cli.Helper
	Message string `cli:"m,message" usage:"Message"`
	Success bool   `cli:"s,success" usage:"Success status"`
}

type rebootT struct {
	cli.Helper
	Serial string `cli:"s,serial" usage:"Device serial number"`
}

type rebootAllT struct {
	cli.Helper
}

var title string = fmt.Sprintf(`
meraki-reboot %s
  Reboot shit tons of Meraki devices because Meraki is terrible`, Version)

var rootCmd = &cli.Command{
	Desc: title,
	Argv: func() interface{} { return new(argT) },
	Fn: func(ctx *cli.Context) error {
		return nil
	},
}

var devicesCmd = &cli.Command{
	Name: "devices",
	Desc: "List all matched devices",
	Argv: func() interface{} { return new(devicesT) },
	Fn: func(ctx *cli.Context) error {
		args := ctx.RootArgv().(*argT)
		c := ctx.Color()
		exclusions := util.SplitRemoveEmpty(args.Exclusions, ",")
		orgID, err := meraki.GetOrganizationID(args.OrgName)
		if err != nil {
			return err
		}
		networkID, err := meraki.GetNetworkID(orgID, args.NetworkName)
		if err != nil {
			return err
		}
		devices, err := meraki.GetNetworkDevices(networkID, exclusions)
		if err != nil {
			return err
		}
		for _, device := range devices {
			ctx.String(fmt.Sprintf("%s: %s\n", c.Green(c.Bold(device.Serial)), c.Grey(device.Name)))
		}
		return nil
	},
}

var rebootCmd = &cli.Command{
	Name: "reboot",
	Desc: "Reboot one device",
	Argv: func() interface{} { return new(rebootT) },
	Fn: func(ctx *cli.Context) error {
		args := ctx.Argv().(*rebootT)
		c := ctx.Color()
		device := meraki.GetDevice(args.Serial)
		success, err := meraki.RebootDevice(args.Serial)
		if err != nil {
			return fmt.Errorf("Failed to reboot device %s (%s)\n%s", c.Red(c.Bold(device.Name)), c.Grey(args.Serial), err)
		}
		if !success {
			return fmt.Errorf("Failed to reboot device %s (%s)\n", c.Red(c.Bold(device.Name)), c.Grey(args.Serial))
		}
		ctx.String("Rebooted device %s (%s)\n", c.Green(c.Bold(device.Name)), c.Grey(args.Serial))
		return nil
	},
}

var slackTestCmd = &cli.Command{
	Name: "slack",
	Desc: "Send a test Slack message",
	Argv: func() interface{} { return new(slackT) },
	Fn: func(ctx *cli.Context) error {
		args := ctx.Argv().(*slackT)
		c := ctx.Color()
		title := "Test Message from Meraki Reboot"
		err := SendWebhook(SlackArgs{Org: "Test", Title: &title, Msg: args.Message, Success: &args.Success})
		if err == nil {
			ctx.String(c.Green(c.Bold("Slack message sent")))
			return err
		}
		return err
	},
}

var rebootAllCmd = &cli.Command{
	Name: "reboot-all",
	Desc: "Reboot all devices",
	Argv: func() interface{} { return new(rebootAllT) },
	Fn: func(ctx *cli.Context) error {
		args := ctx.RootArgv().(*argT)
		c := ctx.Color()
		exclusions := util.SplitRemoveEmpty(args.Exclusions, ",")
		orgID, err := meraki.GetOrganizationID(args.OrgName)
		if err != nil {
			log.Error(err)
			return err
		}
		networkID, err := meraki.GetNetworkID(orgID, args.NetworkName)
		if err != nil {
			log.Error(err)
			return err
		}
		devices, err := meraki.GetNetworkDevices(networkID, exclusions)
		if err != nil {
			log.Error(err)
			return err
		}
		count := 0
		for _, device := range devices {
			success, err := meraki.RebootDevice(device.Serial)
			if err != nil {
				msg := fmt.Sprintf("Failed to reboot device %s (%s)\n%s", device.Name, device.Serial, err)
				log.Error(msg)

				SendWebhook(SlackArgs{Org: args.OrgName, Msg: msg, Success: &success})
				return fmt.Errorf("Failed to reboot device %s (%s)\n%s", c.Red(c.Bold(device.Name)), c.Grey(device.Serial), err)
			}
			if !success {
				msg := fmt.Sprintf("Failed to reboot device %s (%s)", device.Name, device.Serial)
				log.Error(msg)
				SendWebhook(SlackArgs{Org: args.OrgName, Msg: msg, Success: &success})
				return fmt.Errorf("Failed to reboot device %s (%s)\n", c.Red(c.Bold(device.Name)), c.Grey(device.Serial))
			}
			count += 1
			log.Info("Rebooted device %s (%s)", device.Name, device.Serial)
			ctx.String("Rebooted device %s (%s)\n", c.Green(c.Bold(device.Name)), c.Grey(device.Serial))
		}

		msg := fmt.Sprintf(
			"Rebooted %d devices for organization '%s', network '%s'",
			count,
			args.OrgName,
			args.NetworkName,
		)
		log.Info(msg)
		success := true
		SendWebhook(SlackArgs{Org: args.OrgName, Msg: msg, Success: &success})
		ctx.String(
			"Rebooted %d devices for organization '%s', network '%s'\n",
			c.Green(c.Bold(count)),
			c.Blue(c.Bold(args.OrgName)),
			c.Yellow(c.Bold(args.NetworkName)),
		)
		return nil
	},
}
