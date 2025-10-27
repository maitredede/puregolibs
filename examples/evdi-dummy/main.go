package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/maitredede/puregolibs/evdi"
	edid2 "github.com/srlehn/edid"
)

func main() {
	flag.Parse()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	edidList := [][]byte{
		// evdi.DummyEdid[:],
		evdi.EDIDv1_1280x800,
		// evdi.EDIDv1_1440x900,
		// evdi.EDIDv1_1600x900,
		// evdi.EDIDv1_1680x1050,
		// evdi.EDIDv2_1280x720,
		// evdi.EDIDv2_1920x1080,
		// evdi.EDIDv2_3840x2160,
	}
	for _, bin := range edidList {
		edid, err := edid2.New(bin)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		slog.Info(fmt.Sprintf("%+v", edid))
	}

	device, err := evdi.OpenAttachedToNone()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer device.Close()

	// if err := device.EnableCursorEvents(true); err != nil {
	// 	slog.Error(err.Error())
	// 	os.Exit(1)
	// }

	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	// defer stop()
	ctx := context.TODO()

	if err := device.RunDummy(ctx); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
