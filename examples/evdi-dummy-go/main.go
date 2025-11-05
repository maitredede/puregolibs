package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/maitredede/puregolibs/libevdi"
	"github.com/maitredede/puregolibs/resources"
)

func main() {
	flag.Parse()

	slog.SetLogLoggerLevel(slog.LevelDebug)
	// libevdi.SetLogging(slog.Default())

	slog.Info(fmt.Sprintf("uid=%d euid=%d gid=%d egid=%d", os.Getuid(), os.Geteuid(), os.Getgid(), os.Getegid()))

	slog.Info("main: opening evdi device")
	h := libevdi.OpenAttachedToNone()
	if h == nil {
		slog.Error("open failed")
		os.Exit(1)
	}
	defer libevdi.Close(h)

	dw := 1280
	dh := 800
	edid := resources.EDIDv1_1280x800

	libevdi.Connect(h, edid, dw*dh)
	defer libevdi.Disconnect(h)

	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	// defer stop()
	// ctx2, cancel := context.WithTimeout(ctx, 45*time.Second)
	// defer cancel()

	// slog.Info("main: starting dummy")
	// if err := device.RunDummy(ctx2); err != nil {
	// 	slog.Error(err.Error())
	// 	os.Exit(1)
	// }
}
