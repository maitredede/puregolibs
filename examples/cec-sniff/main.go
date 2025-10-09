package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/maitredede/puregolibs/cec"
)

func main() {
	flag.Parse()

	appCallbacks := cec.Callbacks{
		LogMessage: func(cbparam any, message cec.LogMessage) {
			fmt.Printf("log: %v: %s\n", message.Level, message.Message)
		},
		KeyPress: func(cbparam any, key cec.Keypress) {
			fmt.Printf("keypress: %v\n", key)
		},
		CommandReceived: func(cbparam any, command cec.Command) {
			fmt.Printf("cmdReceived: %v\n", command)
		},
		ConfigurationChanged: func(cbparam any, configuration cec.Configuration) {
			fmt.Printf("cfgChanged: %v\n", configuration)
		},
		Alert: func(cbparam any, alert cec.Alert, param cec.Parameter) {
			fmt.Printf("alert: %v %v\n", alert, param)
		},
		MenuStateChanged: func(cbparam any, state cec.MenuState) int32 {
			fmt.Printf("menuStateChanged: %v\n", state)
			return 0
		},
		SourceActivated: func(cbparam any, logicalAddress cec.LogicalAddress, activated bool) {
			fmt.Printf("sourceActivated: %v => %v\n", logicalAddress, activated)
		},
		CommandHandler: func(cbparam any, command cec.Command) int32 {
			fmt.Printf("commandHandler: %v\n", command)
			return 0
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	cecc, err := cec.Open("", "", appCallbacks)
	if err != nil {
		fmt.Printf("main: cec open failed: %v\n", err)
		os.Exit(1)
	}
	defer cecc.Close()

	<-ctx.Done()
}
