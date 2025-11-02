//go:build linux

package libevdi

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/maitredede/puregolibs/drm"
	"github.com/prometheus/procfs"
	"golang.org/x/sys/unix"
)

var (
	evdiInvalidHandle *Handle = nil

	cardUsage []*Handle = make([]*Handle, EvdiUsageLength)
)

func getGenericDevice() int {
	deviceIndex := EvdiInvalidDeviceIndex

	deviceIndex = findUnusedCardFor("")
	if deviceIndex == EvdiInvalidDeviceIndex {
		evdiLogInfo("creating card in /sys/devices/platform")
		writeAddDevice("1")
		deviceIndex = findUnusedCardFor("")
	}
	return deviceIndex
}

func getDeviceAttachedToUsb(sysfsParentDevice string) int {
	panic("TODO")
}

func findUnusedCardFor(parentDevice string) int {
	const evdiPlatformRoot = "/sys/bus/platform/devices"
	deviceIndex := EvdiInvalidDeviceIndex

	files, err := os.ReadDir(evdiPlatformRoot)
	if err != nil {
		panic(err)
	}
	for _, d := range files {
		if !strings.HasPrefix(d.Name(), "evdi") {
			continue
		}
		evdiPath := filepath.Join(evdiPlatformRoot, d.Name())
		if !isCorrectParentDevice(evdiPath, parentDevice) {
			continue
		}
		evdiDrmPath := filepath.Join(evdiPath, "drm")
		devIndex := getDrmDeviceIndex(evdiDrmPath)

		assert(devIndex < EvdiUsageLength && devIndex >= 0)

		if cardUsage[devIndex] == evdiInvalidHandle {
			deviceIndex = devIndex
			break
		}
	}

	return deviceIndex
}

func isCorrectParentDevice(dirname string, parentDevice string) bool {
	linkPath := filepath.Join(dirname, "device")
	if len(parentDevice) == 0 {
		ret := unix.Access(linkPath, unix.F_OK)
		evdiLogDebug("isCorrectParentDevice: access to '%s': %v", linkPath, ret)
		return ret != nil
	}
	linkResolution, err := os.Readlink(linkPath)
	if err != nil {
		evdiLogDebug("readlink %s error %v", linkPath, err)
		return false
	}

	if len(parentDevice) < 2 {
		return false
	}

	_ = linkResolution
	panic("WIP")
}

func getDrmDeviceIndex(evdiSysfsDrmDir string) int {
	devIndex := EvdiInvalidDeviceIndex
	// err := filepath.WalkDir(evdiSysfsDrmDir, func(path string, d fs.DirEntry, err error) error {
	// 	if strings.HasPrefix(d.Name(), "card") {
	// 		num, err := strconv.ParseUint(d.Name()[4:], 10, 32)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		devIndex = int(num)
	// 		return filepath.SkipAll
	// 	}
	// 	return nil
	// })
	// if err != nil {
	// 	evdiLogError(fmt.Sprintf("failed to open dir %s: %v", evdiSysfsDrmDir, err))
	// }
	files, err := os.ReadDir(evdiSysfsDrmDir)
	if err != nil {
		evdiLogError("failed to open dir %s: %v", evdiSysfsDrmDir, err)
		return devIndex
	}
	for _, d := range files {
		name := d.Name()
		if strings.HasPrefix(name, "card") {
			num, err := strconv.ParseUint(name[4:], 10, 32)
			if err != nil {
				panic(err)
			}
			devIndex = int(num)
		}
	}

	return devIndex
}

func writeAddDevice(buffer string) (int, error) {
	addDevices, err := os.OpenFile("/sys/devices/evdi/add", os.O_WRONLY, 0)
	if err != nil {
		evdiLogDebug("writeAddDevice open failed: %v", err)
		return 0, err
	}
	defer addDevices.Close()

	n, err := addDevices.WriteString(buffer)
	if err != nil {
		evdiLogDebug("writeAddDevice write failed: %v", err)
		return 0, err
	}
	return n, nil
}

func openDevice(device int) (*os.File, error) {
	dev := fmt.Sprintf("/dev/dri/card%d", device)
	if XorgRunning() {
		waitForMaster(dev)
	}
	fd, err := waitForDevice(dev)
	if err != nil {
		return nil, err
	}

	// if err := drmIoctl(fd.Fd(), DRM_IOCTL_DROP_MASTER, 0); err == 0 {
	if err := drm.DropMaster(fd); err == nil {
		evdiLogInfo("dropped master on %s", dev)
	}
	return fd, nil
}

func isEvdi(fd *os.File) bool {
	v, err := drm.GetVersion(fd)
	if err != nil {
		evdiLogDebug("isEvdi: drm.GetVersion err=%v", err)
		return false
	}
	return strings.HasPrefix(v.Name, "evdi")
}

func isEvdiCompatible(fd *os.File) bool {

	evdiLogInfo("LibEvdi Version (%d.%d.%d) (go %s)", libEvdiVersionMajor, libEvdiVersionMinor, libEvdiVersionPatch, runtime.Version())

	ver, err := drm.GetVersion(fd)
	if err != nil {
		evdiLogInfo("can't get evdi version: %v", err)
		return false
	}

	evdiLogInfo("Evdi Version (%d.%d.%d)", ver.Major, ver.Minor, ver.Patch)

	if ver.Major == evdiModuleCompatibilityVersionMajor &&
		ver.Minor >= evdiModuleCompatibilityVersionMinor {
		return true
	}

	evdiLogInfo("Doesn't match LibEvdi compatibility one (%d.%d.%d)", evdiModuleCompatibilityVersionMajor, evdiModuleCompatibilityVersionMinor, evdiModuleCompatibilityVersionPatch)

	return false
}

func waitForMaster(devicePath string) {
	totalWaitUS := 5_000_000
	sleepIntervalUs := 100_000
	cnt := totalWaitUS / sleepIntervalUs

	hasMaster := false
	for {
		hasMaster = deviceHasMaster(devicePath)
		cnt--
		if hasMaster || cnt < 0 {
			break
		}
		time.Sleep(time.Duration(sleepIntervalUs) * time.Microsecond)
	}
	if !hasMaster {
		evdiLogInfo("wait for master timed out")
	}
}

func deviceHasMaster(deviceFilePath string) bool {
	myself := os.Getpid()
	result := false

	procs, err := procfs.AllProcs()
	if err != nil {
		panic(err)
	}
	for _, p := range procs {
		if p.PID == myself {
			continue
		}
		if processOpenedFiles(p, deviceFilePath) {
			result = true
			break
		}
		if processOpenedDevice(p, deviceFilePath) {
			result = true
			break
		}
	}
	return result
}

func processOpenedFiles(p procfs.Proc, deviceFilePath string) bool {
	isMatch := false
	targets, err := p.FileDescriptorTargets()
	if err != nil {
		evdiLogDebug("processOpenedFiles %v: %v", p, err)
		return false
	}
	for _, t := range targets {
		if t == deviceFilePath {
			isMatch = true
		}
	}
	return isMatch
}

func processOpenedDevice(p procfs.Proc, deviceFilePath string) bool {
	maps, err := p.ProcMaps()
	if err != nil {
		evdiLogDebug("processOpenedDevice %v: %v", p, err)
		return false
	}
	for _, m := range maps {
		if m.Pathname == deviceFilePath {
			return true
		}
	}
	return false
}

func waitForDevice(devicePath string) (*os.File, error) {
	var f *os.File
	var err error

	totalWaitUS := 5000000
	SleepIntervalUS := 100000
	cnt := totalWaitUS / SleepIntervalUS

	for {
		cnt--
		f, err = openAsSlave(devicePath)
		if f != nil || cnt < 0 {
			break
		}
		time.Sleep(time.Duration(SleepIntervalUS) * time.Microsecond)
	}
	if err != nil {
		evdiLogError("failed to open a device: %v", err)
	}
	return f, err
}

func openAsSlave(devicePath string) (*os.File, error) {
	f, err := os.OpenFile(devicePath, os.O_RDWR|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil, err
	}
	if drm.IsMaster(f) {
		evdiLogInfo("process has master on %s", devicePath)
		// err = drmIoctl(f.Fd(), DRM_IOCTL_DROP_MASTER, 0)
		err = drm.DropMaster(f)
	}
	if err != nil {
		evdiLogInfo("drop master on %s failed, err: %v", devicePath, err)
		f.Close()
		return nil, err
	}
	if drm.IsMaster(f) {
		evdiLogInfo("drop master on %s failed, err: %v", devicePath, err)
		f.Close()
		return nil, err
	}

	evdiLogInfo("opened %s as slave drm device", devicePath)
	return f, nil
}
