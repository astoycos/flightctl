package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/flightctl/flightctl/internal/agent"
	"github.com/flightctl/flightctl/pkg/log"
	"github.com/sirupsen/logrus"
)

func main() {
	command := NewAgentCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

type agentCmd struct {
	log        *logrus.Logger
	config     *agent.Config
	configFile string
}

func NewAgentCommand() *agentCmd {
	a := &agentCmd{
		log:    log.InitLogs(),
		config: agent.NewDefault(),
	}

	flag.StringVar(&a.configFile, "config", agent.DefaultConfigFile, fmt.Sprintf("path to config file: default: %s", agent.DefaultConfigFile))

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Println("This program starts an agent with the specified configuration. Below are the available flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if err := a.config.ParseConfigFile(a.configFile); err != nil {
		a.log.Fatalf("Error parsing config: %v", err)
	}
	if err := a.config.Validate(); err != nil {
		a.log.Fatalf("Error validating config: %v", err)
	}

	logLvl, err := logrus.ParseLevel(a.config.LogLevel)
	if err != nil {
		logLvl = logrus.InfoLevel
	}
	a.log.SetLevel(logLvl)

	return a
}

func (a *agentCmd) Execute() error {
	agentInstance := agent.New(a.log, a.config)
	if err := agentInstance.Run(context.Background()); err != nil {
		a.log.Fatalf("running device agent: %v", err)
	}
	return nil
}
