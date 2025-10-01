package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/maitredede/puregolibs/libnfc"
)

var (
	verbose   bool
	intrusive bool
)

func main() {
	flag.BoolVar(&verbose, "v", false, "set verbose display")
	flag.BoolVar(&intrusive, "i", false, "allow intrusive scan")
	flag.Parse()

	if intrusive {
		// This has to be done before the call to nfc_init()
		if err := os.Setenv("LIBNFC_INTRUSIVE_SCAN", "yes"); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	context, err := libnfc.InitContext()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer context.Close()

	// Display libnfc version
	libnfcVersion := libnfc.Version()
	fmt.Printf("%s uses libnfc %s\n", os.Args[0], libnfcVersion)

	devices, err := context.ListDevices()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(devices) == 0 {
		fmt.Println("no NFC device found.")
		os.Exit(1)
	}
	fmt.Printf("%d NFC device(s) found:\n", len(devices))
	for _, devConString := range devices {
		pnd, err := context.OpenDevice(devConString)
		if err != nil {
			fmt.Printf("nfc_open failed for %s: %v\n", devConString, err)
			continue
		}
		name, err := pnd.Name()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		connStr, err := pnd.ConnString()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("- %s:\n    %s\n", name, connStr)

		if verbose {
			infos, err := pnd.GetInformationAbout()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Print(infos)
		}
		if err := pnd.Close(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
