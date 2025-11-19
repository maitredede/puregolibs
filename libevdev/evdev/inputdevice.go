package evdev

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"syscall"

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

// Open creates a new InputDevice from the given path. The input device is
// opened with flag O_RDWR. Returns an error if the device node could not
// be opened or its properties failed to read.
func Open(path string) (*InputDevice, error) {
	return OpenWithFlags(path, os.O_RDWR)
}

// OpenByNameWithFlags creates a new InputDevice from the device name as reported
// by the kernel. The input device is opened with the specified flags (O_RDONLY etc.).
// It is the responsibility of the user to provide sane flags and handle potential errors
// resulting from inappropriate flag combinations or permissions.
// Returns an error if the name does not exist, or the device node could
// not be opened or its properties failed to read.
func OpenByNameWithFlags(name string, flags int) (*InputDevice, error) {
	devices, err := ListDevicePaths()
	if err != nil {
		return nil, err
	}
	for _, d := range devices {
		if d.Name == name {
			return OpenWithFlags(d.Path, flags)
		}
	}
	return nil, fmt.Errorf("could not find input device with name %q", name)
}

// OpenByName creates a new InputDevice from the device name as reported by the kernel.
// The input device is opened with flag O_RDWR.
// Returns an error if the name does not exist, or the device node could
// not be opened or its properties failed to read.
func OpenByName(name string) (*InputDevice, error) {
	return OpenByNameWithFlags(name, os.O_RDWR)
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

// NonBlock sets file descriptor into nonblocking mode.
// This way it is possible to interrupt ReadOne call by closing the device.
// Note: file.Fd() call will set file descriptor back to blocking mode so make sure your program
// is not using any other method than ReadOne after NonBlock call.
func (d *InputDevice) NonBlock() error {
	return syscall.SetNonblock(int(d.file.Fd()), true)
}

// Read and returns a slice of InputEvents from the device.
// It blocks until events has been received or an error has occurred.
func (d *InputDevice) ReadSlice(eventSlice int) ([]InputEvent, error) {
	buffer := make([]byte, eventsize*eventSlice)

	bytesRead, err := d.file.Read(buffer)
	if err != nil {
		return nil, err
	}

	// Calculate how many complete events we actually got
	count := bytesRead / eventsize
	if count == 0 {
		return nil, nil // no complete event in this read
	}

	// Create events slice dynamically
	events := make([]InputEvent, count)

	reader := bytes.NewReader(buffer[:bytesRead])
	if err = binary.Read(reader, binary.LittleEndian, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// ReadOne reads one InputEvent from the device. It blocks until an event has
// been received or an error has occurred.
func (d *InputDevice) ReadOne() (*InputEvent, error) {
	event := InputEvent{}

	err := binary.Read(d.file, binary.LittleEndian, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// WriteOne writes one InputEvent to the device.
// Useful for controlling LEDs of the device
func (d *InputDevice) WriteOne(event *InputEvent) error {
	return binary.Write(d.file, binary.LittleEndian, event)
}

// InputID returns the device's vendor/product/busType/version information as reported by the kernel.
func (d *InputDevice) InputID() (InputID, error) {
	id := InputID{
		BusType: libevdev.GetIDBusType(d.evdev),
		Vendor:  libevdev.GetIDVendor(d.evdev),
		Product: libevdev.GetIDProduct(d.evdev),
		Version: libevdev.GetIDVersion(d.evdev),
	}
	return id, nil
}
