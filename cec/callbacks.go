package cec

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

type nativeICECCallbacks struct {
	logMessage           uintptr // void (CEC_CDECL* logMessage)(void* cbparam, const cec_log_message* message);
	keyPress             uintptr // void (CEC_CDECL* keyPress)(void* cbparam, const cec_keypress* key);
	commandReceived      uintptr // void (CEC_CDECL* commandReceived)(void* cbparam, const cec_command* command);
	configurationChanged uintptr // void (CEC_CDECL* configurationChanged)(void* cbparam, const libcec_configuration* configuration);
	alert                uintptr // void (CEC_CDECL* alert)(void* cbparam, const libcec_alert alert, const libcec_parameter param);
	menuStateChanged     uintptr // int (CEC_CDECL* menuStateChanged)(void* cbparam, const cec_menu_state state);
	sourceActivated      uintptr // void (CEC_CDECL* sourceActivated)(void* cbParam, const cec_logical_address logicalAddress, const uint8_t bActivated);
	commandHandler       uintptr // int (CEC_CDECL* commandHandler)(void* cbparam, const cec_command* command);
}

type LogMessageCallback func(cbparam any, message LogMessage)
type KeyPressCallback func(cbparam any, key Keypress)
type CommandReceivedCallback func(cbparam any, command Command)
type ConfigurationChangedCallback func(cbparam any, configuration Configuration)
type AlertCallback func(cbparam any, alert Alert, param Parameter)
type MenuStateChangedCallback func(cbparam any, state MenuState) int32
type SourceActivatedCallback func(cbparam any, logicalAddress LogicalAddress, activated bool)
type CommandHandlerCallback func(cbparam any, command Command) int32
type Callbacks struct {
	LogMessage           LogMessageCallback
	KeyPress             KeyPressCallback
	CommandReceived      CommandReceivedCallback
	ConfigurationChanged ConfigurationChangedCallback
	Alert                AlertCallback
	MenuStateChanged     MenuStateChangedCallback
	SourceActivated      SourceActivatedCallback
	CommandHandler       CommandHandlerCallback
}

func buildNativeCallbacks(appCallbacks Callbacks) (*nativeICECCallbacks, []func()) {
	cb := new(nativeICECCallbacks)

	logCB, disposeLogMessage := buildCallbackLogMessage(appCallbacks.LogMessage)
	cb.logMessage = logCB

	keyCB, disposeKeyPress := buildCallbackKeyPress(appCallbacks.KeyPress)
	cb.keyPress = keyCB

	cmdCB, disposeCmdReceived := buildCallbackCommandReceived(appCallbacks.CommandReceived)
	cb.commandReceived = cmdCB

	cfgCB, disposeCfgChanged := buildCallbackConfigurationChanged(appCallbacks.ConfigurationChanged)
	cb.configurationChanged = cfgCB

	alertCB, disposeAlert := buildCallbackAlert(appCallbacks.Alert)
	cb.alert = alertCB

	menuCB, disposeMenu := buildCallbackMenuStateChanged(appCallbacks.MenuStateChanged)
	cb.menuStateChanged = menuCB

	sourceCB, disposeSource := buildCallbackSourceActivated(appCallbacks.SourceActivated)
	cb.sourceActivated = sourceCB

	cmdHandlerCB, disposeCmd := buildCallbackCommandHandler(appCallbacks.CommandHandler)
	cb.commandHandler = cmdHandlerCB

	disposables := []func(){
		disposeLogMessage,
		disposeKeyPress,
		disposeCmdReceived,
		disposeCfgChanged,
		disposeAlert,
		disposeMenu,
		disposeSource,
		disposeCmd,
	}

	return cb, disposables
}

func logMessageCallback(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
	logCB := *(*LogMessageCallback)(userData)

	argsSlice := unsafe.Slice(args, 2)
	cbparam := *(*any)(argsSlice[0])
	nativeMsg := *(*nativeLogMessage)(argsSlice[1])

	message := nativeMsg.Go()

	logCB(cbparam, message)

	return 0
}

func buildCallbackLogMessage(logCB LogMessageCallback) (uintptr, func()) {
	if logCB == nil {
		return 0, func() {}
	}

	// allocate the closure function
	var callback unsafe.Pointer
	closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
	if closure == nil {
		panic("closure alloc failed")
	}

	// describe the closure's signature
	var cifCallback ffi.Cif
	/*!
	 * @brief Transfer a log message from libCEC to the client.
	 * @param cbparam             Callback parameter provided when the callbacks were set up
	 * @param message             The message to transfer.
	 */
	// void (CEC_CDECL* logMessage)(void* cbparam, const cec_log_message* message);
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 2, &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}

	// fn will be called, then the closure gets invoked
	fn := ffi.NewCallback(logMessageCallback)

	// prepare the closure
	userData := unsafe.Pointer(&logCB)
	if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, userData, callback); status != ffi.OK {
		panic(status)
	}

	dispose := func() {
		ffi.ClosureFree(closure)
	}

	return uintptr(callback), dispose
}

func keypressCallback(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
	keypressCB := *(*KeyPressCallback)(userData)

	argsSlice := unsafe.Slice(args, 2)
	cbparam := *(*any)(argsSlice[0])
	nativeKey := *(*nativeKeyPress)(argsSlice[1])

	key := nativeKey.Go()

	keypressCB(cbparam, key)

	return 0
}

func buildCallbackKeyPress(keypressCB KeyPressCallback) (uintptr, func()) {
	if keypressCB == nil {
		return 0, func() {}
	}

	// allocate the closure function
	var callback unsafe.Pointer
	closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
	if closure == nil {
		panic("closure alloc failed")
	}

	// describe the closure's signature
	var cifCallback ffi.Cif
	/*!
	 * @brief Transfer a keypress from libCEC to the client.
	 * @param cbparam             Callback parameter provided when the callbacks were set up
	 * @param key                 The keypress to transfer.
	 */
	// void (CEC_CDECL* keyPress)(void* cbparam, const cec_keypress* key);
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 2, &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}

	// fn will be called, then the closure gets invoked
	fn := ffi.NewCallback(keypressCallback)

	// prepare the closure
	userData := unsafe.Pointer(&keypressCB)
	if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, userData, callback); status != ffi.OK {
		panic(status)
	}

	dispose := func() {
		ffi.ClosureFree(closure)
	}

	return uintptr(callback), dispose
}

func commandReceivedCallback(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
	cmdCB := *(*CommandReceivedCallback)(userData)

	argsSlice := unsafe.Slice(args, 2)
	cbparam := *(*any)(argsSlice[0])
	nativeCmd := *(*nativeCommand)(argsSlice[1])

	cmd := nativeCmd.Go()

	cmdCB(cbparam, cmd)

	return 0
}

func buildCallbackCommandReceived(cmdCB CommandReceivedCallback) (uintptr, func()) {
	if cmdCB == nil {
		return 0, func() {}
	}

	// allocate the closure function
	var callback unsafe.Pointer
	closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
	if closure == nil {
		panic("closure alloc failed")
	}
	// describe the closure's signature
	var cifCallback ffi.Cif
	/*!
	 * @brief Transfer a CEC command from libCEC to the client.
	 * @param cbparam             Callback parameter provided when the callbacks were set up
	 * @param command             The command to transfer.
	 */
	// void (CEC_CDECL* commandReceived)(void* cbparam, const cec_command* command);
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 2, &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}

	// fn will be called, then the closure gets invoked
	fn := ffi.NewCallback(commandReceivedCallback)

	// prepare the closure
	userData := unsafe.Pointer(&cmdCB)
	if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, userData, callback); status != ffi.OK {
		panic(status)
	}

	dispose := func() {
		ffi.ClosureFree(closure)
	}

	return uintptr(callback), dispose
}

func configurationChangedCallback(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
	cfgCB := *(*ConfigurationChangedCallback)(userData)

	argsSlice := unsafe.Slice(args, 2)
	cbparam := *(*any)(argsSlice[0])
	nativeCfg := *(*NativeConfiguration)(argsSlice[1])

	cfg := nativeCfg.Go()

	cfgCB(cbparam, cfg)

	return 0
}

func buildCallbackConfigurationChanged(cfgCB ConfigurationChangedCallback) (uintptr, func()) {
	if cfgCB == nil {
		return 0, func() {}
	}

	// allocate the closure function
	var callback unsafe.Pointer
	closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
	if closure == nil {
		panic("closure alloc failed")
	}
	// describe the closure's signature
	var cifCallback ffi.Cif
	/*!
	 * @brief Transfer a changed configuration from libCEC to the client
	 * @param cbparam             Callback parameter provided when the callbacks were set up
	 * @param configuration       The configuration to transfer
	 */
	// void (CEC_CDECL* configurationChanged)(void* cbparam, const libcec_configuration* configuration);
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 2, &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}

	// fn will be called, then the closure gets invoked
	fn := ffi.NewCallback(configurationChangedCallback)

	// prepare the closure
	userData := unsafe.Pointer(&cfgCB)
	if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, userData, callback); status != ffi.OK {
		panic(status)
	}

	dispose := func() {
		ffi.ClosureFree(closure)
	}

	return uintptr(callback), dispose
}

func alertCallback(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
	alertCB := *(*AlertCallback)(userData)

	argsSlice := unsafe.Slice(args, 3)
	cbparam := *(*any)(argsSlice[0])
	alert := *(*Alert)(argsSlice[1])
	//TODO check if pointer dereference is not needed, because arg is "const libcec_parameter param" (struct)
	nativeParam := *(*nativeParameter)(argsSlice[2])

	param := nativeParam.Go()

	alertCB(cbparam, alert, param)

	return 0
}

func buildCallbackAlert(alertCB AlertCallback) (uintptr, func()) {
	if alertCB == nil {
		return 0, func() {}
	}

	// allocate the closure function
	var callback unsafe.Pointer
	closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
	if closure == nil {
		panic("closure alloc failed")
	}

	// describe the closure's signature
	typeCecParameter := ffi.NewType(
		&ffi.TypeSint32,
		&ffi.TypePointer,
	)
	var cifCallback ffi.Cif
	/*!
	 * @brief Transfer a libcec alert message from libCEC to the client
	 * @param cbparam             Callback parameter provided when the callbacks were set up
	 * @param alert               The alert type transfer.
	 * @param data                Misc. additional information.
	 */
	// void (CEC_CDECL* alert)(void* cbparam, const libcec_alert alert, const libcec_parameter param);
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 3, &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &typeCecParameter); status != ffi.OK {
		panic(status)
	}

	// fn will be called, then the closure gets invoked
	fn := ffi.NewCallback(alertCallback)

	// prepare the closure
	userData := unsafe.Pointer(&alertCB)
	if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, userData, callback); status != ffi.OK {
		panic(status)
	}

	dispose := func() {
		ffi.ClosureFree(closure)
	}

	return uintptr(callback), dispose
}

func menuStateChangedCallback(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
	menuCB := *(*MenuStateChangedCallback)(userData)

	argsSlice := unsafe.Slice(args, 2)
	cbparam := *(*any)(argsSlice[0])
	state := *(*MenuState)(argsSlice[1])

	retVal := menuCB(cbparam, state)

	*(*int32)(ret) = retVal

	return 0
}

func buildCallbackMenuStateChanged(menuCB MenuStateChangedCallback) (uintptr, func()) {
	if menuCB == nil {
		return 0, func() {}
	}

	// allocate the closure function
	var callback unsafe.Pointer
	closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
	if closure == nil {
		panic("closure alloc failed")
	}
	// describe the closure's signature
	var cifCallback ffi.Cif
	/*!
	 * @brief Transfer a menu state change to the client.
	 * Transfer a menu state change to the client. If the command returns 1, then the change will be processed by
	 * the busdevice. If 0, then the state of the busdevice won't be changed, and will always be kept 'activated',
	 * @warning CEC does not allow the player to suppress the menu state change on the TV, so the menu on the TV will always be displayed, whatever the return value of this method is.
	 * so keypresses are always routed.
	 * @param cbparam             Callback parameter provided when the callbacks were set up
	 * @param state               The new value.
	 *
	 * @return 1 if libCEC should use this new value, 0 otherwise.
	 */
	// int (CEC_CDECL* menuStateChanged)(void* cbparam, const cec_menu_state state);
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 2, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32); status != ffi.OK {
		panic(status)
	}

	// fn will be called, then the closure gets invoked
	fn := ffi.NewCallback(menuStateChangedCallback)

	// prepare the closure
	userData := unsafe.Pointer(&menuCB)
	if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, userData, callback); status != ffi.OK {
		panic(status)
	}

	dispose := func() {
		ffi.ClosureFree(closure)
	}

	return uintptr(callback), dispose
}

func sourceActivatedCallback(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
	srcCB := *(*SourceActivatedCallback)(userData)

	argsSlice := unsafe.Slice(args, 3)
	cbparam := *(*any)(argsSlice[0])
	logicalAddress := *(*LogicalAddress)(argsSlice[1])
	nativeActivated := *(*int8)(argsSlice[2])

	activated := nativeActivated != 0

	srcCB(cbparam, logicalAddress, activated)

	return 0
}

func buildCallbackSourceActivated(srcCB SourceActivatedCallback) (uintptr, func()) {
	if srcCB == nil {
		return 0, func() {}
	}

	// allocate the closure function
	var callback unsafe.Pointer
	closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
	if closure == nil {
		panic("closure alloc failed")
	}

	// describe the closure's signature
	var cifCallback ffi.Cif
	/*!
	 * @brief Called when a source that's handled by this client is activated.
	 * @param cbparam             Callback parameter provided when the callbacks were set up
	 * @param logicalAddress      The address that was just activated.
	 * @param bActivated          1 if activated, 0 when deactivated.
	 */
	// void (CEC_CDECL* sourceActivated)(void* cbParam, const cec_logical_address logicalAddress, const uint8_t bActivated);
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 3, &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeUint8); status != ffi.OK {
		panic(status)
	}

	// fn will be called, then the closure gets invoked
	fn := ffi.NewCallback(keypressCallback)

	// prepare the closure
	userData := unsafe.Pointer(&srcCB)
	if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, userData, callback); status != ffi.OK {
		panic(status)
	}

	dispose := func() {
		ffi.ClosureFree(closure)
	}

	return uintptr(callback), dispose
}

func commandHandlerCallback(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
	cmdCB := *(*CommandHandlerCallback)(userData)

	argsSlice := unsafe.Slice(args, 2)
	cbparam := *(*any)(argsSlice[0])
	nativeCmd := *(*nativeCommand)(argsSlice[1])

	cmd := nativeCmd.Go()

	retVal := cmdCB(cbparam, cmd)

	*(*int32)(ret) = retVal

	return 0
}

func buildCallbackCommandHandler(cmdCB CommandHandlerCallback) (uintptr, func()) {
	if cmdCB == nil {
		return 0, func() {}
	}

	// allocate the closure function
	var callback unsafe.Pointer
	closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
	if closure == nil {
		panic("closure alloc failed")
	}
	// describe the closure's signature
	var cifCallback ffi.Cif
	/*!
	 * @brief Allow the client handle a CEC command instead of libcec.
	 * @param cbparam             Callback parameter provided when the callbacks were set up
	 * @param command             The command to handle.
	 *
	 * @return 1 if the command has been handled and if libCEC should not take any action
	 */
	// int (CEC_CDECL* commandHandler)(void* cbparam, const cec_command* command);
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 2, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}

	// fn will be called, then the closure gets invoked
	fn := ffi.NewCallback(commandHandlerCallback)

	// prepare the closure
	userData := unsafe.Pointer(&cmdCB)
	if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, userData, callback); status != ffi.OK {
		panic(status)
	}

	dispose := func() {
		ffi.ClosureFree(closure)
	}

	return uintptr(callback), dispose
}
