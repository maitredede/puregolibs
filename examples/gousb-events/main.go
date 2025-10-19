package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/maitredede/puregolibs/gousb"
)

var (
	debug int
)

func main() {
	flag.IntVar(&debug, "debug", 1, "libusb debug level (0..3)")
	flag.Parse()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	slog.Info(fmt.Sprintf("libusb version: %v", gousb.GetVersion()))
	slog.Info(fmt.Sprintf("cap: Hotplug: %v", gousb.HasCapability(gousb.CapHasHotplug)))
	slog.Info(fmt.Sprintf("cap: HID access: %v", gousb.HasCapability(gousb.CapHasHIDAccess)))
	slog.Info(fmt.Sprintf("cap: Kernel detach: %v", gousb.HasCapability(gousb.CapSupportsDetachKernelDriver)))

	// Only one context should be needed for an application.  It should always be closed.
	ctx, err := gousb.Init()
	if err != nil {
		slog.Error(fmt.Sprintf("ctx: %s", err))
		os.Exit(1)
	}
	defer ctx.Close()

	// Debugging can be turned on; this shows some of the inner workings of the libusb package.
	ctx.SetDebug(gousb.LogLevel(debug))

	events := gousb.HotplugEventDeviceArrived | gousb.HotplugEventDeviceLeft
	flags := gousb.HotplugFlagsEnumerate
	hotHandle, err := ctx.RegisterCallback(events, flags, gousb.HotplugMatchAny, gousb.HotplugMatchAny, gousb.HotplugMatchAny, hotplugCallback)
	if err != nil {
		slog.Error(fmt.Sprintf("hotplug: %s", err))
		os.Exit(1)
	}
	defer hotHandle.Deregister()

	c, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	func(c context.Context) {
		for {
			select {
			case <-c.Done():
				return
			default:
			}

			n, err := ctx.HandleEventsTimeout(1 * time.Second)
			if n != 0 || err != nil {
				slog.Info(fmt.Sprintf("events: returned n=%v err=%v", n, err))
			}
		}
	}(c)
}

func hotplugCallback(ctx *gousb.Context, device *gousb.Device, event gousb.HotplugEvent, userData any) {
	slog.Warn(fmt.Sprintf("TODO: hotplugCallback evt=%v dev=%v", event, device))
}
