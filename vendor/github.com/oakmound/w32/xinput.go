package w32

import (
	"errors"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	xinput       *windows.LazyDLL
	xinputEnable *windows.LazyProc
	//xinputGetAudioDeviceIds         *windows.LazyProc
	xinputGetBatteryInformation     *windows.LazyProc
	xinputGetCapabilities           *windows.LazyProc
	xinputGetDSoundAudioDeviceGuids *windows.LazyProc
	xinputGetKeystroke              *windows.LazyProc
	xinputGetState                  *windows.LazyProc
	xinputSetState                  *windows.LazyProc
)

// InitXInput populates xinput procs for later calls.
func InitXInput() error {
	// Todo: this takes a lot of work to handle these different files, it'd be
	// nice if we do it without having to check multiple potential dlls.
	names := []string{"xinput1_4.dll", "xinput1_3.dll", "xinput9_1_0.dll", "xinputuap.dll"}
	for _, n := range names {
		dll, err := tryLoad(n)
		if err == nil {
			xinput = dll
			break
		}
	}
	if xinput == nil {
		return errors.New("Unable to find xinput dll")
	}
	xinputEnable = xinput.NewProc("XInputEnable")
	//xinputGetAudioDeviceIds = xinput.NewProc("XInputGetAudioDeviceIds")
	xinputGetBatteryInformation = xinput.NewProc("XInputGetBatteryInformation")
	xinputGetCapabilities = xinput.NewProc("XInputGetCapabilities")
	xinputGetDSoundAudioDeviceGuids = xinput.NewProc("XInputGetDSoundAudioDeviceGuids")
	xinputGetKeystroke = xinput.NewProc("XInputGetKeystroke")
	xinputGetState = xinput.NewProc("XInputGetState")
	xinputSetState = xinput.NewProc("XInputSetState")
	return nil
}

func tryLoad(dllName string) (*windows.LazyDLL, error) {
	dll := &windows.LazyDLL{
		Name:   dllName,
		System: true,
	}
	return dll, dll.Load()
}

type XInputState struct {
	PacketNumber uint32
	Gamepad      XInputGamepad
}

type XInputGamepad struct {
	Buttons      uint16
	LeftTrigger  uint8
	RightTrigger uint8
	ThumbLX      int16
	ThumbLY      int16
	ThumbRX      int16
	ThumbRY      int16
}

type XInputKeystroke struct {
	VirtualKey uint16
	Unicode    uint16
	Flags      uint16
	UserIndex  uint8
	HidCode    uint8
}

type XInputCapabilities struct {
	Type      uint8
	SubType   uint8
	Flags     uint16
	Gamepad   XInputGamepad
	Vibration XInputVibration
}

type XInputBatteryInformation struct {
	BatteryType  uint8
	BatteryLevel uint8
}

type XInputVibration struct {
	LeftMotorSpeed  uint16
	RightMotorSpeed uint16
}

const (
	// Bitmasks for the joysticks buttons, determines what has
	// been pressed on the joystick, these need to be mapped
	// to whatever device you're using instead of an xbox 360
	// joystick
	XINPUT_GAMEPAD_DPAD_UP        = 0x0001
	XINPUT_GAMEPAD_DPAD_DOWN      = 0x0002
	XINPUT_GAMEPAD_DPAD_LEFT      = 0x0004
	XINPUT_GAMEPAD_DPAD_RIGHT     = 0x0008
	XINPUT_GAMEPAD_START          = 0x0010
	XINPUT_GAMEPAD_BACK           = 0x0020
	XINPUT_GAMEPAD_LEFT_THUMB     = 0x0040
	XINPUT_GAMEPAD_RIGHT_THUMB    = 0x0080
	XINPUT_GAMEPAD_LEFT_SHOULDER  = 0x0100
	XINPUT_GAMEPAD_RIGHT_SHOULDER = 0x0200
	XINPUT_GAMEPAD_A              = 0x1000
	XINPUT_GAMEPAD_B              = 0x2000
	XINPUT_GAMEPAD_X              = 0x4000
	XINPUT_GAMEPAD_Y              = 0x8000

	// Defines the flags used to determine if the user is pushing
	// down on a button, not holding a button, etc
	XINPUT_KEYSTROKE_KEYDOWN = 0x0001
	XINPUT_KEYSTROKE_KEYUP   = 0x0002
	XINPUT_KEYSTROKE_REPEAT  = 0x0004

	// Defines the codes which are returned by XInputGetKeystroke
	VK_PAD_A                = 0x5800
	VK_PAD_B                = 0x5801
	VK_PAD_X                = 0x5802
	VK_PAD_Y                = 0x5803
	VK_PAD_RSHOULDER        = 0x5804
	VK_PAD_LSHOULDER        = 0x5805
	VK_PAD_LTRIGGER         = 0x5806
	VK_PAD_RTRIGGER         = 0x5807
	VK_PAD_DPAD_UP          = 0x5810
	VK_PAD_DPAD_DOWN        = 0x5811
	VK_PAD_DPAD_LEFT        = 0x5812
	VK_PAD_DPAD_RIGHT       = 0x5813
	VK_PAD_START            = 0x5814
	VK_PAD_BACK             = 0x5815
	VK_PAD_LTHUMB_PRESS     = 0x5816
	VK_PAD_RTHUMB_PRESS     = 0x5817
	VK_PAD_LTHUMB_UP        = 0x5820
	VK_PAD_LTHUMB_DOWN      = 0x5821
	VK_PAD_LTHUMB_RIGHT     = 0x5822
	VK_PAD_LTHUMB_LEFT      = 0x5823
	VK_PAD_LTHUMB_UPLEFT    = 0x5824
	VK_PAD_LTHUMB_UPRIGHT   = 0x5825
	VK_PAD_LTHUMB_DOWNRIGHT = 0x5826
	VK_PAD_LTHUMB_DOWNLEFT  = 0x5827
	VK_PAD_RTHUMB_UP        = 0x5830
	VK_PAD_RTHUMB_DOWN      = 0x5831
	VK_PAD_RTHUMB_RIGHT     = 0x5832
	VK_PAD_RTHUMB_LEFT      = 0x5833
	VK_PAD_RTHUMB_UPLEFT    = 0x5834
	VK_PAD_RTHUMB_UPRIGHT   = 0x5835
	VK_PAD_RTHUMB_DOWNRIGHT = 0x5836
	VK_PAD_RTHUMB_DOWNLEFT  = 0x5837

	// Deadzones are for analogue joystick controls on the joypad
	// which determine when input should be assumed to be in the
	// middle of the pad. This is a threshold to stop a joypad
	// controlling the game when the player isn't touching the
	// controls.
	XINPUT_GAMEPAD_LEFT_THUMB_DEADZONE  = 7849
	XINPUT_GAMEPAD_RIGHT_THUMB_DEADZONE = 8689
	XINPUT_GAMEPAD_TRIGGER_THRESHOLD    = 30

	// Defines what type of abilities the type of joystick has
	// DEVTYPE_GAMEPAD is available for all joysticks, however
	// there may be more specific identifiers for other joysticks
	// which are being used.
	XINPUT_DEVTYPE_GAMEPAD         = 0x01
	XINPUT_DEVSUBTYPE_GAMEPAD      = 0x01
	XINPUT_DEVSUBTYPE_WHEEL        = 0x02
	XINPUT_DEVSUBTYPE_ARCADE_STICK = 0x03
	XINPUT_DEVSUBTYPE_FLIGHT_SICK  = 0x04
	XINPUT_DEVSUBTYPE_DANCE_PAD    = 0x05
	XINPUT_DEVSUBTYPE_GUITAR       = 0x06
	XINPUT_DEVSUBTYPE_DRUM_KIT     = 0x08

	// These are used with the XInputGetCapabilities function to
	// determine the abilities to the joystick which has been
	// plugged in.
	XINPUT_CAPS_VOICE_SUPPORTED = 0x0004
	XINPUT_FLAG_GAMEPAD         = 0x00000001

	// Defines the status of the battery if one is used in the
	// attached joystick. The first two define if the joystick
	// supports a battery. Disconnected means that the joystick
	// isn't connected. Wired shows that the joystick is a wired
	// joystick.
	BATTERY_DEVTYPE_GAMEPAD   = 0x00
	BATTERY_DEVTYPE_HEADSET   = 0x01
	BATTERY_TYPE_DISCONNECTED = 0x00
	BATTERY_TYPE_WIRED        = 0x01
	BATTERY_TYPE_ALKALINE     = 0x02
	BATTERY_TYPE_NIMH         = 0x03
	BATTERY_TYPE_UNKNOWN      = 0xFF
	BATTERY_LEVEL_EMPTY       = 0x00
	BATTERY_LEVEL_LOW         = 0x01
	BATTERY_LEVEL_MEDIUM      = 0x02
	BATTERY_LEVEL_FULL        = 0x03

	// How many joysticks can be used with this library. Games that
	// use the xinput library will not go over this number.
	XUSER_MAX_COUNT = 4
	XUSER_INDEX_ANY = 0x000000FF
)

func XInputGetDSoundAudioDeviceGuids(userIndex uint32, render, capture *GUID) error {
	r, _, _ := xinputGetDSoundAudioDeviceGuids.Call(
		uintptr(userIndex),
		uintptr(unsafe.Pointer(render)),
		uintptr(unsafe.Pointer(capture)),
	)
	if r == 0 {
		return nil
	}
	return syscall.Errno(r)
}

func XInputGetKeystroke(userIndex uint32, keystroke *XInputKeystroke) error {
	r, _, _ := xinputGetKeystroke.Call(
		uintptr(userIndex),
		uintptr(0),
		uintptr(unsafe.Pointer(keystroke)),
	)
	if r == 0 {
		return nil
	}
	return syscall.Errno(r)
}

func XInputGetCapabilities(userIndex, flags uint32, capabilities *XInputCapabilities) error {
	r, _, _ := xinputGetCapabilities.Call(
		uintptr(userIndex),
		uintptr(flags),
		uintptr(unsafe.Pointer(capabilities)),
	)
	if r == 0 {
		return nil
	}
	return syscall.Errno(r)
}

func XInputGetBatteryInformation(userIndex uint32, devType uint8, battInfo *XInputBatteryInformation) error {
	r, _, _ := xinputGetBatteryInformation.Call(
		uintptr(userIndex),
		uintptr(devType),
		uintptr(unsafe.Pointer(battInfo)),
	)
	if r == 0 {
		return nil
	}
	return syscall.Errno(r)
}

func XInputEnable(enable bool) error {
	var enableInt uintptr
	if enable {
		enableInt++
	}
	r, _, _ := xinputEnable.Call(uintptr(enableInt))
	if r == 0 {
		return nil
	}
	return syscall.Errno(r)
}

func XInputGetState(userIndex uint32, state *XInputState) error {
	r, _, _ := xinputGetState.Call(
		uintptr(userIndex),
		uintptr(unsafe.Pointer(state)),
	)
	if r == 0 {
		return nil
	}
	return syscall.Errno(r)
}

func XInputSetState(userIndex uint32, vibration *XInputVibration) error {
	r, _, _ := xinputSetState.Call(
		uintptr(userIndex),
		uintptr(unsafe.Pointer(vibration)),
	)
	if r == 0 {
		return nil
	}
	return syscall.Errno(r)
}
