package evdi

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/maitredede/puregolibs/resources"
	"github.com/neilotoole/slogt"
)

func TestPyEvdiDummyMonitor(t *testing.T) {
	initLib()

	log := slogt.New(t)
	slog.SetDefault(log)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	SetLogging(func(s string) {
		log.Info(s)
	})

	options := Options{
		Headless:         true,
		ResolutionWidth:  1280,
		ResolutionHeight: 800,
		RefreshRate:      60,
		// EDID:             EDIDv2_1920x1080,
		EDID: resources.EDIDv1_1280x800,
		// EDID:     nil,
		FPSLimit: 60,
	}

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
	defer cancel()

	RunMain(ctx, options)
}
