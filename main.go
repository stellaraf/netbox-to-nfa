package main

import (
	"fmt"
	"os"

	"github.com/mkideal/cli"
	log "github.com/sirupsen/logrus"
	"stellar.af/netbox-to-nfa/util"
)

func logSetup(fileName string) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	util.Check("Error accessing log file '%s'", err, fileName)
	log.SetOutput(file)
}

func init() {
	log.Debug("Checking environment variables...")
	util.CheckEnv("NETBOX_URL", true)
	util.CheckEnv("NETBOX_TOKEN", true)
	util.CheckEnv("NETBOX_NFA_ROLE", true)
	util.CheckEnv("NFA_URL", true)
	util.CheckEnv("NFA_USERNAME", true)
	util.CheckEnv("NFA_PASSWORD", true)
	util.CheckEnv("NB2NFA_EXCLUDED_RANGES", true)
	logFile := util.GetEnv("NB2NFA_LOGFILE")
	if logFile == "" {
		logFile = "/var/log/nb2nfa.log"
	}
	logSetup(logFile)
	log.Debug("All required environment variables are present")
	log.Info("Starting netbox-to-nfa...")
}

func main() {
	if err := cli.Root(rootCmd,
		cli.Tree(cli.HelpCommand("display help information")),
		cli.Tree(purgeCmd),
		cli.Tree(syncCmd),
		cli.Tree(prefixesCmd),
		cli.Tree(filtersCmd),
		cli.Tree(cfgCmd),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
