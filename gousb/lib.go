package gousb

import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

var (
	initLckOnce sync.Mutex
	initPtr     uintptr
	initError   error
)

func getSystemLibrary() string {
	switch runtime.GOOS {
	// case "darwin":
	// 	return "libusb.dylib"
	case "linux":
		return "libusb-1.0.so"
	case "windows":
		return "libusb-1.0.dll"
	default:
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}
}

func mustGetSymbol(sym string) uintptr {
	ptr, err := getSymbol(sym)
	if err != nil {
		panic(err)
	}
	return ptr
}

func libInitFuncs() {
	purego.RegisterLibFunc(&libusbInit, initPtr, "libusb_init")
	if sym, err := getSymbol("libusb_init_context"); err == nil {
		purego.RegisterFunc(&libusbInitContext, sym)
	} else {
		libusbInitContext = func(ctx *libusbContext, options *NativeLibusbInitOption, numOptions int32) int32 {
			if options != nil || numOptions > 0 {
				fmt.Println("WARN : libusb_init_context not found (old version of libusb). falling back to libusb_init")
			}
			return libusbInit(ctx)
		}
	}
	purego.RegisterLibFunc(&libusbExit, initPtr, "libusb_exit")
	purego.RegisterLibFunc(&libusbSetDebug, initPtr, "libusb_set_debug")
	purego.RegisterLibFunc(&libusbSetLogCb, initPtr, "libusb_set_log_cb")
	purego.RegisterLibFunc(&libusbGetVersion, initPtr, "libusb_get_version")
	purego.RegisterLibFunc(&libusbHasCapability, initPtr, "libusb_has_capability")
	purego.RegisterLibFunc(&libusbErrorName, initPtr, "libusb_error_name")
	purego.RegisterLibFunc(&libusbSetLocale, initPtr, "libusb_setlocale")
	purego.RegisterLibFunc(&libusbStrError, initPtr, "libusb_strerror")

	purego.RegisterLibFunc(&libusbGetDeviceList, initPtr, "libusb_get_device_list")
	purego.RegisterLibFunc(&libusbFreeDeviceList, initPtr, "libusb_free_device_list")
	purego.RegisterLibFunc(&libusbRefDevice, initPtr, "libusb_ref_device")
	purego.RegisterLibFunc(&libusbUnrefDevice, initPtr, "libusb_unref_device")

	// libusb_get_device_string
	purego.RegisterLibFunc(&libusbGetConfiguration, initPtr, "libusb_get_configuration")
	purego.RegisterLibFunc(&libusbGetDeviceDescriptor, initPtr, "libusb_get_device_descriptor")
	purego.RegisterLibFunc(&libusbGetActiveConfigDescriptor, initPtr, "libusb_get_active_config_descriptor")
	purego.RegisterLibFunc(&libusbGetConfigDescriptor, initPtr, "libusb_get_config_descriptor")
	purego.RegisterLibFunc(&libusbGetConfigDescriptorByValue, initPtr, "libusb_get_config_descriptor_by_value")
	purego.RegisterLibFunc(&libusbFreeConfigDescriptor, initPtr, "libusb_free_config_descriptor")
	// libusb_get_ss_endpoint_companion_descriptor
	// libusb_free_ss_endpoint_companion_descriptor
	// libusb_get_bos_descriptor
	// libusb_free_bos_descriptor
	// libusb_get_usb_2_0_extension_descriptor
	// libusb_free_usb_2_0_extension_descriptor
	// libusb_get_ss_usb_device_capability_descriptor
	// libusb_free_ss_usb_device_capability_descriptor
	// libusb_get_ssplus_usb_device_capability_descriptor
	// libusb_free_ssplus_usb_device_capability_descriptor
	// libusb_get_container_id_descriptor
	// libusb_free_container_id_descriptor
	// libusb_get_platform_descriptor
	// libusb_free_platform_descriptor
	purego.RegisterLibFunc(&libusbGetBusNumber, initPtr, "libusb_get_bus_number")
	purego.RegisterLibFunc(&libusbGetPortNumber, initPtr, "libusb_get_port_number")
	purego.RegisterLibFunc(&libusbGetPortNumbers, initPtr, "libusb_get_port_numbers")
	purego.RegisterLibFunc(&libusbGetPortPath, initPtr, "libusb_get_port_path")
	purego.RegisterLibFunc(&libusbGetParent, initPtr, "libusb_get_parent")
	purego.RegisterLibFunc(&libusbGetDeviceAddress, initPtr, "libusb_get_device_address")
	purego.RegisterLibFunc(&libusbGetDeviceSpeed, initPtr, "libusb_get_device_speed")
	purego.RegisterLibFunc(&libusbGetMaxPacketSize, initPtr, "libusb_get_max_packet_size")
	purego.RegisterLibFunc(&libusbGetMaxIsoPacketSize, initPtr, "libusb_get_max_iso_packet_size")
	if sym, err := getSymbol("libusb_get_max_alt_packet_size"); err == nil {
		purego.RegisterFunc(&libusbGetMaxAltPacketSize, sym)
	} else {
		libusbGetMaxAltPacketSize = func(dev libusbDevice, interfaceNumber int32, alternateSettings int32, endpoint uint8) int32 {
			panic("libusb_get_max_alt_packet_size not found (old version of libusb).")
		}
	}

	// libusb_get_interface_association_descriptors
	// libusb_get_active_interface_association_descriptors
	// libusb_free_interface_association_descriptors

	// libusb_wrap_sys_device
	purego.RegisterLibFunc(&libusbOpen, initPtr, "libusb_open")
	purego.RegisterLibFunc(&libusbClose, initPtr, "libusb_close")
	purego.RegisterLibFunc(&libusbGetDevice, initPtr, "libusb_get_device")

	purego.RegisterLibFunc(&libusbSetConfiguration, initPtr, "libusb_set_configuration")
	purego.RegisterLibFunc(&libusbClaimInterface, initPtr, "libusb_claim_interface")
	purego.RegisterLibFunc(&libusbReleaseInterface, initPtr, "libusb_release_interface")

	purego.RegisterLibFunc(&libusbOpenDeviceWithVidPid, initPtr, "libusb_open_device_with_vid_pid")

	purego.RegisterLibFunc(&libusbSetInterfaceAltSetting, initPtr, "libusb_set_interface_alt_setting")
	purego.RegisterLibFunc(&libusbClearHalt, initPtr, "libusb_clear_halt")
	purego.RegisterLibFunc(&libusbResetDevice, initPtr, "libusb_reset_device")

	purego.RegisterLibFunc(&libusbAllocStreams, initPtr, "libusb_alloc_streams")
	purego.RegisterLibFunc(&libusbFreeStreams, initPtr, "libusb_free_streams")

	// libusb_dev_mem_alloc
	// libusb_dev_mem_free

	// libusb_kernel_driver_active
	// libusb_detach_kernel_driver
	// libusb_attach_kernel_driver
	// libusb_set_auto_detach_kernel_driver

	// libusb_alloc_transfer
	// libusb_submit_transfer
	// libusb_cancel_transfer
	// libusb_free_transfer
	// libusb_transfer_set_stream_id
	// libusb_transfer_get_stream_id

	// libusb_control_transfer
	// libusb_bulk_transfer
	// libusb_interrupt_transfer

	// libusb_get_string_descriptor_ascii

	// libusb_try_lock_events
	// libusb_lock_events
	// libusb_unlock_events
	// libusb_event_handling_ok
	// libusb_event_handler_active
	// libusb_interrupt_event_handler
	// libusb_lock_event_waiters
	// libusb_unlock_event_waiters
	// libusb_wait_for_event

	// libusb_handle_events_timeout
	purego.RegisterLibFunc(&libusbHandleEventsTimeoutCompleted, initPtr, "libusb_handle_events_timeout_completed")
	// libusb_handle_events
	// libusb_handle_events_completed
	// libusb_handle_events_locked
	// libusb_pollfds_handle_timeouts
	// libusb_get_next_timeout

	// libusb_get_pollfds
	// libusb_free_pollfds
	// libusb_set_pollfd_notifiers

	purego.RegisterLibFunc(&libusbHotplugRegisterCallback, initPtr, "libusb_hotplug_register_callback")
	purego.RegisterLibFunc(&libusbHotplugDeregisterCallback, initPtr, "libusb_hotplug_deregister_callback")
	purego.RegisterLibFunc(&libusbHotplugGetUserData, initPtr, "libusb_hotplug_get_user_data")

	purego.RegisterLibFunc(&libusbSetOption, initPtr, "libusb_set_option")
	purego.RegisterLibFunc(&libusbSetOptionInt32, initPtr, "libusb_set_option")
	purego.RegisterLibFunc(&libusbSetOptionPtr, initPtr, "libusb_set_option")
}

var (
	libusbInit          func(ctx *libusbContext) int32
	libusbInitContext   func(ctx *libusbContext, options *NativeLibusbInitOption, numOptions int32) int32
	libusbExit          func(ctx libusbContext)
	libusbSetDebug      func(ctx libusbContext, level LogLevel)
	libusbSetLogCb      func(ctx libusbContext, cb unsafe.Pointer, mode int32)
	libusbGetVersion    func() *libusbVersion
	libusbHasCapability func(cap Capability) bool
	libusbErrorName     func(errorCode int32) string
	libusbSetLocale     func(ctx libusbContext, locale string) int32
	libusbStrError      func(errorCode int32) string

	libusbGetDeviceList  func(ctx libusbContext, list **libusbDevice) int32
	libusbFreeDeviceList func(list *libusbDevice, unrefDevices int32)
	libusbRefDevice      func(device libusbDevice) libusbDevice
	libusbUnrefDevice    func(device libusbDevice)

	// libusb_get_device_string
	libusbGetConfiguration           func(device libusbDeviceHandle, config *int32) int32
	libusbGetDeviceDescriptor        func(dev libusbDevice, desc *libusbDeviceDescriptor) int32
	libusbGetActiveConfigDescriptor  func(dev libusbDevice, config **libusbConfigDescriptor) int32
	libusbGetConfigDescriptor        func(dev libusbDevice, configIndex uint8, config **libusbConfigDescriptor) int32
	libusbGetConfigDescriptorByValue func(dev libusbDevice, configurationValue uint8, config **libusbConfigDescriptor) int32
	libusbFreeConfigDescriptor       func(config *libusbConfigDescriptor)
	// libusb_get_ss_endpoint_companion_descriptor
	// libusb_free_ss_endpoint_companion_descriptor
	// libusb_get_bos_descriptor
	// libusb_free_bos_descriptor
	// libusb_get_usb_2_0_extension_descriptor
	// libusb_free_usb_2_0_extension_descriptor
	// libusb_get_ss_usb_device_capability_descriptor
	// libusb_free_ss_usb_device_capability_descriptor
	// libusb_get_ssplus_usb_device_capability_descriptor
	// libusb_free_ssplus_usb_device_capability_descriptor
	// libusb_get_container_id_descriptor
	// libusb_free_container_id_descriptor
	// libusb_get_platform_descriptor
	// libusb_free_platform_descriptor
	libusbGetBusNumber        func(dev libusbDevice) uint8
	libusbGetPortNumber       func(dev libusbDevice) uint8
	libusbGetPortNumbers      func(dev libusbDevice, portNumbers *uint8, portNumbersLen int32) int32
	libusbGetPortPath         func(ctx libusbContext, dev libusbDevice, path *uint8, pathLength uint8) int32
	libusbGetParent           func(dev libusbDevice) libusbDevice
	libusbGetDeviceAddress    func(dev libusbDevice) uint8
	libusbGetDeviceSpeed      func(dev libusbDevice) int32
	libusbGetMaxPacketSize    func(dev libusbDevice, endpoint uint8) int32
	libusbGetMaxIsoPacketSize func(dev libusbDevice, endpoint uint8) int32
	libusbGetMaxAltPacketSize func(dev libusbDevice, interfaceNumber int32, alternateSettings int32, endpoint uint8) int32

	// libusb_get_interface_association_descriptors
	// libusb_get_active_interface_association_descriptors
	// libusb_free_interface_association_descriptors

	// libusb_wrap_sys_device
	libusbOpen      func(device libusbDevice, deviceHandle *libusbDeviceHandle) int32
	libusbClose     func(deviceHandle libusbDeviceHandle)
	libusbGetDevice func(deviceHandle libusbDeviceHandle) libusbDevice

	libusbSetConfiguration func(devHandle libusbDeviceHandle, configuration int32) int32
	libusbClaimInterface   func(devHandle libusbDeviceHandle, interfaceNumber int32) int32
	libusbReleaseInterface func(devHandle libusbDeviceHandle, interfaceNumber int32) int32

	libusbOpenDeviceWithVidPid func(ctx libusbContext, vendorID uint16, productID uint16) libusbDeviceHandle

	libusbSetInterfaceAltSetting func(devHandle libusbDeviceHandle, interfaceNumber int32, alternateSettings int32) int32
	libusbClearHalt              func(devHandle libusbDeviceHandle, endpoint byte) int32
	libusbResetDevice            func(devHandle libusbDeviceHandle) int32

	libusbAllocStreams func(devHandle libusbDeviceHandle, numStreams uint32, endpoints *byte, numEnpoints int32) int32
	libusbFreeStreams  func(devHandle libusbDeviceHandle, endpoints *byte, numEnpoints int32) int32

	// libusb_dev_mem_alloc
	// libusb_dev_mem_free

	// libusb_kernel_driver_active
	// libusb_detach_kernel_driver
	// libusb_attach_kernel_driver
	// libusb_set_auto_detach_kernel_driver

	// libusb_alloc_transfer
	// libusb_submit_transfer
	// libusb_cancel_transfer
	// libusb_free_transfer
	// libusb_transfer_set_stream_id
	// libusb_transfer_get_stream_id

	// libusb_control_transfer
	// libusb_bulk_transfer
	// libusb_interrupt_transfer

	// libusb_get_string_descriptor_ascii

	// libusb_try_lock_events
	// libusb_lock_events
	// libusb_unlock_events
	// libusb_event_handling_ok
	// libusb_event_handler_active
	// libusb_interrupt_event_handler
	// libusb_lock_event_waiters
	// libusb_unlock_event_waiters
	// libusb_wait_for_event

	// libusb_handle_events_timeout
	libusbHandleEventsTimeoutCompleted func(ctx libusbContext, tv *NativeTimeval, completed *int32) int32
	// libusb_handle_events
	// libusb_handle_events_completed
	// libusb_handle_events_locked
	// libusb_pollfds_handle_timeouts
	// libusb_get_next_timeout

	// libusb_get_pollfds
	// libusb_free_pollfds
	// libusb_set_pollfd_notifiers

	libusbHotplugRegisterCallback   func(ctx libusbContext, events HotplugEvent, flags HotplugFlags, vendorId, productID, devClass int32, cbfn unsafe.Pointer, userData unsafe.Pointer, h *hotplugCallbackHandle) int32
	libusbHotplugDeregisterCallback func(ctx libusbContext, handle hotplugCallbackHandle)
	libusbHotplugGetUserData        func(ctx libusbContext, handle hotplugCallbackHandle) unsafe.Pointer

	libusbSetOption      func(ctx libusbContext, option libusbOption, values ...any) int32
	libusbSetOptionInt32 func(ctx libusbContext, option libusbOption, value int32) int32
	libusbSetOptionPtr   func(ctx libusbContext, option libusbOption, value unsafe.Pointer) int32
)
