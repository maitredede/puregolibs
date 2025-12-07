package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/maitredede/puregolibs/vlc"
)

func main() {
	flag.Parse()

	version := vlc.GetVersion()
	slog.Info(fmt.Sprintf("vlc version: %s", version))

	instance, err := vlc.New(nil)
	if err != nil {
		slog.Error(fmt.Sprintf("new vlc error: %v", err))
		os.Exit(1)
	}
	defer instance.Close()

	m := instance.NewMediaFromLocation("http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4")
	mp := m.NewMediaPlayer()
	m.Close()

	mp.Play()

	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	ctx, cancel := context.WithTimeout(sigCtx, 10*time.Second)
	defer cancel()

	<-ctx.Done()

	mp.Stop()
	mp.Close()
}
