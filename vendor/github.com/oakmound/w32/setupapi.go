package w32

import (
	"reflect"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type HDEVINFO HANDLE

type SP_DEVINFO_DATA struct {
	CbSize    uint32
	ClassGuid GUID
	DevInst   uint32
	Reserved  *uint
}

type SP_DEVICE_INTERFACE_DATA struct {
	CbSize    uint32
	ClassGuid GUID
	DevInst   uint32
	Reserved  *uint
}

type SP_DEVICE_INTERFACE_DETAIL_DATA struct {
	CbSize     uint32
	DevicePath string
}

type SPDeviceInformationData struct {
	SP_DEVINFO_DATA
	devInfo HDEVINFO
}

// SetupDiClassGuidsFromNameEx retrieves the GUIDs associated with the specified class name. This resulting list contains the classes currently installed on a local or remote computer.
func SetupDiClassGuidsFromNameEx(className string, machineName string) ([]GUID, error) {
	requiredSize := uint32(0)
	err := setupDiClassGuidsFromNameEx(className, nil, 0, &requiredSize, machineName, 0)

	rets := make([]GUID, requiredSize, requiredSize)
	err = setupDiClassGuidsFromNameEx(className, &rets[0], 1, &requiredSize, machineName, 0)
	return rets, err
}

// SetupDiEnumDeviceInfo returns a SP_DEVINFO_DATA structure that specifies a device information element in a device information set.
func SetupDiEnumDeviceInfo(di HDEVINFO, memberIndex uint32) (*SP_DEVINFO_DATA, error) {
	did := new(SP_DEVINFO_DATA)

	did.CbSize = uint32(unsafe.Sizeof(did))

	err := setupDiEnumDeviceInfo(HANDLE(di), memberIndex, did)
	return did, err
}

// InstanceID retrieves the device instance ID that is associated with a device information element
func (did *SPDeviceInformationData) InstanceID() (string, error) {
	requiredSize := uint32(0)
	err := setupDiGetDeviceInstanceId(HANDLE(did.devInfo), &did.SP_DEVINFO_DATA, nil, 0, &requiredSize)

	buff := make([]uint16, requiredSize)
	err = setupDiGetDeviceInstanceId(HANDLE(did.devInfo), &did.SP_DEVINFO_DATA, unsafe.Pointer(&buff[0]), uint32(len(buff)), &requiredSize)
	if err != nil {
		return "", err
	}

	return windows.UTF16ToString(buff[:]), err
}

// SetupDiGetClassDevsEx returns a handle to a device information set that contains requested device information elements for a local or a remote computer.
func SetupDiGetClassDevsEx(ClassGuid GUID, Enumerator string, hwndParent uintptr, Flags uint32, DeviceInfoSet uintptr, MachineName string, reserved uint32) (HDEVINFO, error) {
	enumerator := &Enumerator

	if Enumerator == "" {
		enumerator = nil
	}

	hDevInfo, err := setupDiGetClassDevsEx(&ClassGuid, enumerator, hwndParent, uint32(Flags), DeviceInfoSet, MachineName, 0)
	return HDEVINFO(hDevInfo), err
}

var (
	modsetupapi = syscall.NewLazyDLL("setupapi.dll")

	procSetupDiClassGuidsFromNameExW     = modsetupapi.NewProc("SetupDiClassGuidsFromNameExW")
	procSetupDiGetClassDevsExW           = modsetupapi.NewProc("SetupDiGetClassDevsExW")
	procSetupDiEnumDeviceInfo            = modsetupapi.NewProc("SetupDiEnumDeviceInfo")
	procSetupDiEnumDeviceInterfaces      = modsetupapi.NewProc("SetupDiEnumDeviceInterfaces")
	procSetupDiGetDeviceInterfaceDetail  = modsetupapi.NewProc("SetupDiGetDeviceInterfaceDetail")
	procSetupDiGetDeviceInstanceIdW      = modsetupapi.NewProc("SetupDiGetDeviceInstanceIdW")
	procSetupDiGetDeviceRegistryProperty = modsetupapi.NewProc("SetupDiGetDeviceRegistryProperty")
)

func SetupDiGetDeviceRegistryProperty(DeviceInfoSet HDEVINFO, DeviceInfoData *SP_DEVINFO_DATA, Property uint32, PropertyRegDataType *uint32, PropertyBuffer []byte, PropertyBufferSize uint32, RequiredSize *uint32) bool {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&PropertyBuffer))
	ret, _, _ := procSetupDiGetDeviceRegistryProperty.Call(
		uintptr(DeviceInfoSet),
		uintptr(unsafe.Pointer(DeviceInfoData)),
		uintptr(Property),
		uintptr(unsafe.Pointer(PropertyRegDataType)),
		uintptr(unsafe.Pointer(hdr.Data)),
		uintptr(PropertyBufferSize),
		uintptr(unsafe.Pointer(RequiredSize)),
	)
	return ret != 0
}

func setupDiClassGuidsFromNameEx(ClassName string, guid *GUID, size uint32, required_size *uint32, machineName string, reserved uint32) (err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(ClassName)
	if err != nil {
		return
	}
	var _p1 *uint16
	_p1, err = syscall.UTF16PtrFromString(machineName)
	if err != nil {
		return
	}
	return _setupDiClassGuidsFromNameEx(_p0, guid, size, required_size, _p1, reserved)
}

func _setupDiClassGuidsFromNameEx(ClassName *uint16, guid *GUID, size uint32, required_size *uint32, machineName *uint16, reserved uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procSetupDiClassGuidsFromNameExW.Addr(), 6, uintptr(unsafe.Pointer(ClassName)), uintptr(unsafe.Pointer(guid)), uintptr(size), uintptr(unsafe.Pointer(required_size)), uintptr(unsafe.Pointer(machineName)), uintptr(reserved))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func setupDiGetClassDevsEx(ClassGuid *GUID, Enumerator *string, hwndParent uintptr, Flags uint32, DeviceInfoSet uintptr, MachineName string, reserved uint32) (handle HANDLE, err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(MachineName)
	if err != nil {
		return
	}
	return _setupDiGetClassDevsEx(ClassGuid, Enumerator, hwndParent, Flags, DeviceInfoSet, _p0, reserved)
}

func _setupDiGetClassDevsEx(ClassGuid *GUID, Enumerator *string, hwndParent uintptr, Flags uint32, DeviceInfoSet uintptr, MachineName *uint16, reserved uint32) (handle HANDLE, err error) {
	r0, _, e1 := syscall.Syscall9(procSetupDiGetClassDevsExW.Addr(), 7, uintptr(unsafe.Pointer(ClassGuid)), uintptr(unsafe.Pointer(Enumerator)), uintptr(hwndParent), uintptr(Flags), uintptr(DeviceInfoSet), uintptr(unsafe.Pointer(MachineName)), uintptr(reserved), 0, 0)
	handle = HANDLE(r0)
	if handle == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func setupDiEnumDeviceInfo(DeviceInfoSet HANDLE, MemberIndex uint32, DeviceInfoData *SP_DEVINFO_DATA) (err error) {
	r1, _, e1 := syscall.Syscall(procSetupDiEnumDeviceInfo.Addr(), 3, uintptr(DeviceInfoSet), uintptr(MemberIndex), uintptr(unsafe.Pointer(DeviceInfoData)))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetupDiEnumDeviceInterfaces(DeviceInfoSet HDEVINFO, DeviceInfoData *SP_DEVINFO_DATA, InterfaceClassGuid *GUID, MemberIndex uint32, DeviceInterfaceData *SP_DEVICE_INTERFACE_DATA) bool {
	ret, _, _ := procSetupDiEnumDeviceInterfaces.Call(
		uintptr(DeviceInfoSet),
		uintptr(unsafe.Pointer(DeviceInfoData)),
		uintptr(unsafe.Pointer(InterfaceClassGuid)),
		uintptr(MemberIndex),
		uintptr(unsafe.Pointer(DeviceInterfaceData)))
	return ret != 0
}

func SetupDiGetDeviceInterfaceDetail(DeviceInfoSet HDEVINFO, DeviceInterfaceData *SP_DEVICE_INTERFACE_DATA, DeviceInterfaceDetailData *SP_DEVICE_INTERFACE_DETAIL_DATA, DeviceInterfaceDetailDataSize uint32, RequiredSize *uint32, DeviceInfoData *SP_DEVINFO_DATA) bool {
	ret, _, _ := procSetupDiGetDeviceInterfaceDetail.Call(
		uintptr(DeviceInfoSet),
		uintptr(unsafe.Pointer(DeviceInterfaceData)),
		uintptr(unsafe.Pointer(DeviceInterfaceDetailData)),
		uintptr(DeviceInterfaceDetailDataSize),
		uintptr(unsafe.Pointer(RequiredSize)),
		uintptr(unsafe.Pointer(DeviceInfoData)),
	)
	return ret != 0
}

func setupDiGetDeviceInstanceId(DeviceInfoSet HANDLE, DeviceInfoData *SP_DEVINFO_DATA, DeviceInstanceId unsafe.Pointer, DeviceInstanceIdSize uint32, RequiredSize *uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procSetupDiGetDeviceInstanceIdW.Addr(), 5, uintptr(DeviceInfoSet), uintptr(unsafe.Pointer(DeviceInfoData)), uintptr(DeviceInstanceId), uintptr(DeviceInstanceIdSize), uintptr(unsafe.Pointer(RequiredSize)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
