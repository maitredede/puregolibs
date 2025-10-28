package evdi

import (
	"context"
	"log/slog"
)

func GetAvailableEvdiCard() int32 {
	for i := int32(0); i < 20; i++ {
		r := libEvdiCheckDevice(i)
		if r == DeviceStatusAvailable {
			return i
		}
	}
	libEvdiAddDevice()
	for i := int32(0); i < 20; i++ {
		r := libEvdiCheckDevice(i)
		if r == DeviceStatusAvailable {
			return i
		}
	}
	return -1
}

type Options struct {
	Headless         bool
	ResolutionWidth  int
	ResolutionHeight int
	RefreshRate      int
	EDID             []byte
	FPSLimit         int
}

var DefaultOptions Options = Options{
	Headless:         false,
	ResolutionWidth:  1920,
	ResolutionHeight: 1080,
	RefreshRate:      60,
	EDID:             nil,
	FPSLimit:         60,
}

func RunMain(ctx context.Context, options Options) {
	device := GetAvailableEvdiCard()
	card, err := NewCard(device)
	if err != nil {
		panic(err)
	}
	area := uint32(options.ResolutionWidth * options.ResolutionHeight)
	/*connectRet :=*/ card.Connect(options.EDID, area, area*uint32(options.RefreshRate))

	slog.Info("running headless")
	for {
		run := true
		select {
		case <-ctx.Done():
			run = false
		default:
		}
		if !run {
			break
		}
		card.HandleEvents(100)
	}

	card.Disconnect()
	card.Close()
}
