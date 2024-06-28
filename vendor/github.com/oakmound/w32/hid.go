package w32

import "syscall"

var (
	hiddll = syscall.NewLazyDLL("hid.dll")

	procHidD_GetAttributes         = hiddll.NewProc("HidD_GetAttributes")
	procHidD_GetSerialNumberString = hiddll.NewProc("HidD_GetSerialNumberString")
	procHidD_GetManufacturerString = hiddll.NewProc("HidD_GetManufacturerString")
	procHidD_GetProductString      = hiddll.NewProc("HidD_GetProductString")
	procHidD_SetFeature            = hiddll.NewProc("HidD_SetFeature")
	procHidD_GetFeature            = hiddll.NewProc("HidD_GetFeature")
	procHidD_GetIndexedString      = hiddll.NewProc("HidD_GetIndexedString")
	procHidD_GetPreparsedData      = hiddll.NewProc("HidD_GetPreparsedData")
	procHidD_FreePreparsedData     = hiddll.NewProc("HidD_FreePreparsedData")
	procHidP_GetCaps               = hiddll.NewProc("HidP_GetCaps")
	procHidD_SetNumInputBuffers    = hiddll.NewProc("HidD_SetNumInputBuffers")
)
