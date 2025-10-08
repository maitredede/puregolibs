package cec

// Device structure
type Device struct {
	OSDName         string
	Vendor          string
	LogicalAddress  int
	ActiveSource    bool
	CECVersion      string
	PowerStatus     string
	PhysicalAddress string
	Language        string
}

var logicalNames = []string{"TV", "Recording", "Recording2", "Tuner",
	"Playback", "Audio", "Tuner2", "Tuner3",
	"Playback2", "Recording3", "Tuner4", "Playback3",
	"Reserved", "Reserved2", "Free", "Broadcast"}

var vendorList = map[uint64]string{0x000039: "Toshiba", 0x0000F0: "Samsung",
	0x0005CD: "Denon", 0x000678: "Marantz", 0x000982: "Loewe", 0x0009B0: "Onkyo",
	0x000CB8: "Medion", 0x000CE7: "Toshiba", 0x001582: "Pulse Eight",
	0x0020C7: "Akai", 0x002467: "Aoc", 0x008045: "Panasonic", 0x00903E: "Philips",
	0x009053: "Daewoo", 0x00A0DE: "Yamaha", 0x00D0D5: "Grundig",
	0x00E036: "Pioneer", 0x00E091: "LG", 0x08001F: "Sharp", 0x080046: "Sony",
	0x18C086: "Broadcom", 0x6B746D: "Vizio", 0x8065E9: "Benq",
	0x9C645E: "Harman Kardon"}

var keyList = map[int]string{0x00: "Select", 0x01: "Up", 0x02: "Down", 0x03: "Left",
	0x04: "Right", 0x05: "RightUp", 0x06: "RightDown", 0x07: "LeftUp",
	0x08: "LeftDown", 0x09: "RootMenu", 0x0A: "SetupMenu", 0x0B: "ContentsMenu",
	0x0C: "FavoriteMenu", 0x0D: "Exit", 0x20: "0", 0x21: "1", 0x22: "2", 0x23: "3",
	0x24: "4", 0x25: "5", 0x26: "6", 0x27: "7", 0x28: "8", 0x29: "9", 0x2A: "Dot",
	0x2B: "Enter", 0x2C: "Clear", 0x2F: "NextFavorite", 0x30: "ChannelUp",
	0x31: "ChannelDown", 0x32: "PreviousChannel", 0x33: "SoundSelect",
	0x34: "InputSelect", 0x35: "DisplayInformation", 0x36: "Help",
	0x37: "PageUp", 0x38: "PageDown", 0x40: "Power", 0x41: "VolumeUp",
	0x42: "VolumeDown", 0x43: "Mute", 0x44: "Play", 0x45: "Stop", 0x46: "Pause",
	0x47: "Record", 0x48: "Rewind", 0x49: "FastForward", 0x4A: "Eject",
	0x4B: "Forward", 0x4C: "Backward", 0x4D: "StopRecord", 0x4E: "PauseRecord",
	0x50: "Angle", 0x51: "SubPicture", 0x52: "VideoOnDemand",
	0x53: "ElectronicProgramGuide", 0x54: "TimerProgramming",
	0x55: "InitialConfiguration", 0x60: "PlayFunction", 0x61: "PausePlay",
	0x62: "RecordFunction", 0x63: "PauseRecordFunction",
	0x64: "StopFunction", 0x65: "Mute",
	0x66: "RestoreVolume", 0x67: "Tune", 0x68: "SelectMedia",
	0x69: "SelectAvInput", 0x6A: "SelectAudioInput", 0x6B: "PowerToggle",
	0x6C: "PowerOff", 0x6D: "PowerOn", 0x71: "Blue", 0x72: "Red", 0x73: "Green",
	0x74: "Yellow", 0x75: "F5", 0x76: "Data", 0x91: "AnReturn",
	0x96: "Max"}
