package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/maitredede/puregolibs/libfreefare"
	"github.com/maitredede/puregolibs/libnfc"
)

var (
	yes     bool
	outfile string
)

func main() {
	flag.BoolVar(&yes, "y", false, "do not ask for confirmation")
	flag.StringVar(&outfile, "o", "", "extrant NDEF message if available in FILE")
	flag.Parse()

	var ndefOut io.Writer
	var messageOut io.StringWriter
	if len(outfile) == 0 || outfile == "-" {
		ndefOut = os.Stdout
		messageOut = os.Stderr
	} else {
		f, err := os.Create(outfile)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		defer f.Close()
		ndefOut = f
		messageOut = os.Stdout
	}

	context, err := libnfc.InitContext()
	if err != nil {
		messageOut.WriteString(err.Error())
		os.Exit(1)
	}
	defer context.Close()

	devices, err := context.ListDevices()
	if err != nil {
		messageOut.WriteString(err.Error())
		os.Exit(1)
	}

	for _, connStr := range devices {
		handleDevice(context, connStr, messageOut, ndefOut)
	}
}

func handleDevice(context *libnfc.NfcContext, connStr string, messageOut io.StringWriter, ndefOut io.Writer) {
	device, err := context.OpenDevice(connStr)
	if err != nil {
		messageOut.WriteString(fmt.Sprintf("device '%s' open failed: %v\n", connStr, err))
		//continue
		return
	}
	defer device.Close()

	tags, err := libfreefare.GetTags(device)
	if err != nil {
		messageOut.WriteString(fmt.Sprintf("error listing MIFARE classic tag: %v\n", err))
		//continue
		return
	}
	defer tags.Close()

	//for _, tag := range tags {
	for i := 0; i < tags.Len(); i++ {
		tag := tags.Get(i)
		typ := tag.TagType()
		if typ != libfreefare.TypeMifareClassic1K && typ != libfreefare.TypeMifareClassic4K {
			continue
		}
		messageOut.WriteString(fmt.Sprintf("found %s with uid %s\n", tag.Name(), tag.UID()))

		readNdef := true

		if readNdef {
			err := readTagNDEF(tag, messageOut, ndefOut)
			if err != nil {
				messageOut.WriteString(err.Error() + "\n")
				continue
			}
		}
	}
}

func readTagNDEF(tag libfreefare.MifareTag, messageOut io.StringWriter, ndefOut io.Writer) error {
	// NFCForum card has a MAD, load it.
	if err := tag.MifareClassicConnect(); err != nil {
		return fmt.Errorf("tag connect error: %w", err)
	}
	defer tag.MifareClassicDisconnect()

	mad, err := tag.ReadMad()
	if err != nil {
		return fmt.Errorf("tag MAD error: %w", err)
	}
	_ = mad
	panic("WIP")
}
