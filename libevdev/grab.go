package libevdev

type GrabMode int32

const (
	GrabModeGrab   GrabMode = 3 // Grab the device if not currently grabbed
	GrabModeUngrab GrabMode = 4 // Ungrab the device if currently grabbed
)

func Grab(evdev Evdev, grab bool) int32 {
	initLib()

	return libevdevGrab(evdev, grab)
}
