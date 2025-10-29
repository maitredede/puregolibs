package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	evdi "github.com/maitredede/puregolibs/evdi/libevdi"
)

func main() {
	flag.Parse()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	slog.Info("main: opening evdi device")
	device, err := evdi.OpenAttachedToNone()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer device.Close()

	slog.Info("main: enabling cursor events")
	device.EnableCursorEvents(true)

	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	// defer stop()
	ctx := context.TODO()

	slog.Info("main: starting dummy")
	if err := device.RunDummy(ctx); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
