package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/maitredede/puregolibs/libevdi"
)

func main() {
	flag.Parse()

	slog.SetLogLoggerLevel(slog.LevelDebug)
	libevdi.SetLogging(slog.Default())

	slog.Info("main: opening evdi device")
	device, err := libevdi.OpenAttachedToNone()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer device.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	ctx2, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	slog.Info("main: starting dummy")
	if err := device.RunDummy(ctx2); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
