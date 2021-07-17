package main

import (
	"fmt"

	"os"

	"github.com/mkideal/cli"
	log "github.com/sirupsen/logrus"
	"stellar.af/meraki-reboot/util"
)

var Version string = "0.2.2"

func logSetup(fileName string) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	util.Check("Error accessing log file '%s'", err, fileName)
	log.SetOutput(file)
}

func init() {
	log.Debug("Checking environment variables...")
	util.CheckEnv("MERAKI_API_KEY", true)
	logFile := util.GetEnv("MERAKI_REBOOT_LOG_FILE")
	if logFile == "" {
		logFile = "/var/log/meraki-reboot.log"
	}
	logSetup(logFile)
	log.Debug("All required environment variables are present")
}

func main() {
	if err := cli.Root(rootCmd,
		cli.Tree(cli.HelpCommand("Display Help Information")),
		cli.Tree(devicesCmd),
		cli.Tree(rebootCmd),
		cli.Tree(rebootAllCmd),
		cli.Tree(slackTestCmd),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
