package libnfc

type Property int16

const (
	NPTimeoutCommand Property = iota // Default command processing timeout
	NPTimeoutATR                     // Timeout between ATR_REQ and ATR_RES
	NPTimeoutCom                     // Timeout value to give up reception from the target in case of no answer
	NPHandleCRC
	NPHandleParity
	NPActivateField
	NPActivateCrypt01
	NPInfiniteSelect
	NPAcceptInvalidFrames
	NPAcceptMultipleFrames
	NPAutoIso14443_4
	NPEasyFraming
	NPForceISO14443A
	NPForceISO14443B
	NPForceSpeed106
)
