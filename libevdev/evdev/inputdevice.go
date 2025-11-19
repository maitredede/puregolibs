package evdev

import (
	"os"

	"github.com/maitredede/puregolibs/libevdev"
)

type InputDevice struct {
	file          *os.File
	evdev         libevdev.Evdev
	driverVersion int32
}

// Close releases the resources held by an InputDevice. After calling this
// function, the InputDevice is no longer operational.
func (d *InputDevice) Close() error {
	libevdev.Free(d.evdev)
	return d.file.Close()
}

// OpenWithFlags creates a new InputDevice from the given path. The input device
// is opened with the specified flags (O_RDONLY etc.).
// It is the responsibility of the user to provide sane flags and handle potential errors
// resulting from inappropriate flag combinations or permissions.
// Returns an error if the device node could not be opened or its properties failed to read.
func OpenWithFlags(path string, flags int) (*InputDevice, error) {
	f, err := os.OpenFile(path, flags, 0)
	if err != nil {
		return nil, err
	}

	evdev := libevdev.New()
	libevdev.SetFd(evdev, f.Fd())

	version := libevdev.GetDriverVersion(evdev)

	d := &InputDevice{
		file:          f,
		evdev:         evdev,
		driverVersion: version,
	}
	return d, nil
}

// Name returns the device's name as reported by the kernel.
func (d *InputDevice) Name() string {
	return libevdev.GetName(d.evdev)
}

// Path returns the device's node path it was opened under.
func (d *InputDevice) Path() string {
	return d.file.Name()
}
