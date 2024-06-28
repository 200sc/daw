// Copyright 2010-2012 The W32 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package w32

import (
	"syscall"
	"unsafe"
)

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procGetModuleHandle    = modkernel32.NewProc("GetModuleHandleW")
	procMulDiv             = modkernel32.NewProc("MulDiv")
	procGetConsoleWindow   = modkernel32.NewProc("GetConsoleWindow")
	procGetCurrentThread   = modkernel32.NewProc("GetCurrentThread")
	procGetLogicalDrives   = modkernel32.NewProc("GetLogicalDrives")
	procGetUserDefaultLCID = modkernel32.NewProc("GetUserDefaultLCID")
	procLstrlen            = modkernel32.NewProc("lstrlenW")
	procLstrcpy            = modkernel32.NewProc("lstrcpyW")
	procGlobalAlloc        = modkernel32.NewProc("GlobalAlloc")
	procGlobalFree         = modkernel32.NewProc("GlobalFree")
	procGlobalLock         = modkernel32.NewProc("GlobalLock")
	procGlobalUnlock       = modkernel32.NewProc("GlobalUnlock")
	procMoveMemory         = modkernel32.NewProc("RtlMoveMemory")
	procFindResource       = modkernel32.NewProc("FindResourceW")
	procSizeofResource     = modkernel32.NewProc("SizeofResource")
	procLockResource       = modkernel32.NewProc("LockResource")
	procLoadResource       = modkernel32.NewProc("LoadResource")
	procGetLastError       = modkernel32.NewProc("GetLastError")
	// procOpenProcess                = modkernel32.NewProc("OpenProcess")
	// procTerminateProcess           = modkernel32.NewProc("TerminateProcess")
	procCloseHandle                = modkernel32.NewProc("CloseHandle")
	procCreateToolhelp32Snapshot   = modkernel32.NewProc("CreateToolhelp32Snapshot")
	procModule32First              = modkernel32.NewProc("Module32FirstW")
	procModule32Next               = modkernel32.NewProc("Module32NextW")
	procGetSystemTimes             = modkernel32.NewProc("GetSystemTimes")
	procGetConsoleScreenBufferInfo = modkernel32.NewProc("GetConsoleScreenBufferInfo")
	procSetConsoleTextAttribute    = modkernel32.NewProc("SetConsoleTextAttribute")
	procGetDiskFreeSpaceEx         = modkernel32.NewProc("GetDiskFreeSpaceExW")
	procGetProcessTimes            = modkernel32.NewProc("GetProcessTimes")
	procSetSystemTime              = modkernel32.NewProc("SetSystemTime")
	procGetSystemTime              = modkernel32.NewProc("GetSystemTime")
	procVirtualAllocEx             = modkernel32.NewProc("VirtualAllocEx")
	procVirtualFreeEx              = modkernel32.NewProc("VirtualFreeEx")
	procWriteProcessMemory         = modkernel32.NewProc("WriteProcessMemory")
	procReadProcessMemory          = modkernel32.NewProc("ReadProcessMemory")
	procQueryPerformanceCounter    = modkernel32.NewProc("QueryPerformanceCounter")
	procQueryPerformanceFrequency  = modkernel32.NewProc("QueryPerformanceFrequency")
	procCreateEvent                = modkernel32.NewProc("CreateEvent")
	procFormatMessage              = modkernel32.NewProc("FormatMessage")
	procCreateFile                 = modkernel32.NewProc("CreateFile")
)

const (
	GENERIC_ALL     = 0x10000000
	GENERIC_EXECUTE = 0x20000000
	GENERIC_WRITE   = 0x40000000
	GENERIC_READ    = 0x80000000

	FILE_SHARE_DELETE = 0x00000004
	FILE_SHARE_READ   = 0x00000001
	FILE_SHARE_WRITE  = 0x00000002

	CREATE_ALWAYS     = 2
	CREATE_NEW        = 1
	OPEN_ALWAYS       = 4
	OPEN_EXISTING     = 3
	TRUNCATE_EXISTING = 5

	FILE_ATTRIBUTE_ARCHIVE   = 0x20
	FILE_ATTRIBUTE_ENCRYPTED = 0x4000
	FILE_ATTRIBUTE_HIDDEN    = 0x2
	FILE_ATTRIBUTE_NORMAL    = 0x80
	FILE_ATTRIBUTE_OFFLINE   = 0x1000
	FILE_ATTRIBUTE_READONLY  = 0x1
	FILE_ATTRIBUTE_SYSTEM    = 0x4
	FILE_ATTRIBUTE_TEMPORARY = 0x100

	FILE_FLAG_BACKUP_SEMANTICS   = 0x02000000
	FILE_FLAG_DELETE_ON_CLOSE    = 0x04000000
	FILE_FLAG_NO_BUFFERING       = 0x20000000
	FILE_FLAG_OPEN_NO_RECALL     = 0x00100000
	FILE_FLAG_OPEN_REPARSE_POINT = 0x00200000
	FILE_FLAG_OVERLAPPED         = 0x40000000
	FILE_FLAG_POSIX_SEMANTICS    = 0x0100000
	FILE_FLAG_RANDOM_ACCESS      = 0x10000000
	FILE_FLAG_SESSION_AWARE      = 0x00800000
	FILE_FLAG_SEQUENTIAL_SCAN    = 0x08000000
	FILE_FLAG_WRITE_THROUGH      = 0x80000000
)

func CreateFile(filename string, desiredAccess, shareMode uint32, security *SECURITY_ATTRIBUTES, creationDisposition, flags uint32, templateFile HANDLE) HANDLE {
	var fn uintptr
	if filename != "" {
		fn = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(filename)))
	}
	ret, _, _ := procCreateFile.Call(fn, uintptr(desiredAccess), uintptr(shareMode), uintptr(unsafe.Pointer(security)), uintptr(creationDisposition), uintptr(flags), uintptr(templateFile))
	return HANDLE(ret)
}

const (
	FORMAT_MESSAGE_ALLOCATE_BUFFER = 0x00000100
	FORMAT_MESSAGE_ARGUMENT_ARRAY  = 0x00002000
	FORMAT_MESSAGE_FROM_HMODULE    = 0x00000800
	FORMAT_MESSAGE_FROM_STRING     = 0x00000400
	FORMAT_MESSAGE_FROM_SYSTEM     = 0x00001000
	FORMAT_MESSAGE_IGNORE_INSERTS  = 0x00000200
	FORMAT_MESSAGE_MAX_WIDTH_MASK  = 0x000000FF
)

func FormatMessage(flags uint32, source uintptr, messageId, languageId uint32, buffer string, size uint32) int {
	var bf uintptr
	if buffer != "" {
		bf = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(buffer)))
	}
	ret, _, _ := procFormatMessage.Call(uintptr(flags), source, uintptr(messageId), uintptr(languageId), bf, uintptr(size), 0)
	return int(ret)
}

func CreateEvent(attributes *SECURITY_ATTRIBUTES, manualReset, initialState bool, name string) HANDLE {
	var mr uintptr
	if manualReset {
		mr = 1
	}
	var is uintptr
	if initialState {
		is = 1
	}
	var nm uintptr
	if name != "" {
		nm = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(name)))
	}
	ret, _, _ := procCreateEvent.Call(uintptr(unsafe.Pointer(attributes)), mr, is, nm)
	return HANDLE(ret)
}

func GetModuleHandle(modulename string) HINSTANCE {
	var mn uintptr
	if modulename != "" {
		mn = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(modulename)))
	}
	ret, _, _ := procGetModuleHandle.Call(mn)
	return HINSTANCE(ret)
}

func MulDiv(number, numerator, denominator int) int {
	ret, _, _ := procMulDiv.Call(
		uintptr(number),
		uintptr(numerator),
		uintptr(denominator))

	return int(ret)
}

func GetConsoleWindow() HWND {
	ret, _, _ := procGetConsoleWindow.Call()

	return HWND(ret)
}

func GetCurrentThread() HANDLE {
	ret, _, _ := procGetCurrentThread.Call()

	return HANDLE(ret)
}

func GetLogicalDrives() uint32 {
	ret, _, _ := procGetLogicalDrives.Call()

	return uint32(ret)
}

func GetUserDefaultLCID() uint32 {
	ret, _, _ := procGetUserDefaultLCID.Call()

	return uint32(ret)
}

func Lstrlen(lpString *uint16) int {
	ret, _, _ := procLstrlen.Call(uintptr(unsafe.Pointer(lpString)))

	return int(ret)
}

func Lstrcpy(buf []uint16, lpString *uint16) {
	procLstrcpy.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(lpString)))
}

func GlobalAlloc(uFlags uint, dwBytes uint32) HGLOBAL {
	ret, _, _ := procGlobalAlloc.Call(
		uintptr(uFlags),
		uintptr(dwBytes))

	if ret == 0 {
		panic("GlobalAlloc failed")
	}

	return HGLOBAL(ret)
}

func GlobalFree(hMem HGLOBAL) {
	ret, _, _ := procGlobalFree.Call(uintptr(hMem))

	if ret != 0 {
		panic("GlobalFree failed")
	}
}

func GlobalLock(hMem HGLOBAL) unsafe.Pointer {
	ret, _, _ := procGlobalLock.Call(uintptr(hMem))

	if ret == 0 {
		panic("GlobalLock failed")
	}

	return unsafe.Pointer(ret)
}

func GlobalUnlock(hMem HGLOBAL) bool {
	ret, _, _ := procGlobalUnlock.Call(uintptr(hMem))

	return ret != 0
}

func MoveMemory(destination, source unsafe.Pointer, length uint32) {
	procMoveMemory.Call(
		uintptr(unsafe.Pointer(destination)),
		uintptr(source),
		uintptr(length))
}

func FindResource(hModule HMODULE, lpName, lpType *uint16) (HRSRC, error) {
	ret, _, _ := procFindResource.Call(
		uintptr(hModule),
		uintptr(unsafe.Pointer(lpName)),
		uintptr(unsafe.Pointer(lpType)))

	if ret == 0 {
		return 0, syscall.GetLastError()
	}

	return HRSRC(ret), nil
}

func SizeofResource(hModule HMODULE, hResInfo HRSRC) uint32 {
	ret, _, _ := procSizeofResource.Call(
		uintptr(hModule),
		uintptr(hResInfo))

	if ret == 0 {
		panic("SizeofResource failed")
	}

	return uint32(ret)
}

func LockResource(hResData HGLOBAL) unsafe.Pointer {
	ret, _, _ := procLockResource.Call(uintptr(hResData))

	if ret == 0 {
		panic("LockResource failed")
	}

	return unsafe.Pointer(ret)
}

func LoadResource(hModule HMODULE, hResInfo HRSRC) HGLOBAL {
	ret, _, _ := procLoadResource.Call(
		uintptr(hModule),
		uintptr(hResInfo))

	if ret == 0 {
		panic("LoadResource failed")
	}

	return HGLOBAL(ret)
}

func GetLastError() uint32 {
	ret, _, _ := procGetLastError.Call()
	return uint32(ret)
}

// func OpenProcess(desiredAccess uint32, inheritHandle bool, processId uint32) HANDLE {
// 	inherit := 0
// 	if inheritHandle {
// 		inherit = 1
// 	}

// 	ret, _, _ := procOpenProcess.Call(
// 		uintptr(desiredAccess),
// 		uintptr(inherit),
// 		uintptr(processId))
// 	return HANDLE(ret)
// }

// func TerminateProcess(hProcess HANDLE, uExitCode uint) bool {
// 	ret, _, _ := procTerminateProcess.Call(
// 		uintptr(hProcess),
// 		uintptr(uExitCode))
// 	return ret != 0
// }

func CloseHandle(object HANDLE) bool {
	ret, _, _ := procCloseHandle.Call(
		uintptr(object))
	return ret != 0
}

func CreateToolhelp32Snapshot(flags, processId uint32) HANDLE {
	ret, _, _ := procCreateToolhelp32Snapshot.Call(
		uintptr(flags),
		uintptr(processId))

	if ret <= 0 {
		return HANDLE(0)
	}

	return HANDLE(ret)
}

func Module32First(snapshot HANDLE, me *MODULEENTRY32) bool {
	ret, _, _ := procModule32First.Call(
		uintptr(snapshot),
		uintptr(unsafe.Pointer(me)))

	return ret != 0
}

func Module32Next(snapshot HANDLE, me *MODULEENTRY32) bool {
	ret, _, _ := procModule32Next.Call(
		uintptr(snapshot),
		uintptr(unsafe.Pointer(me)))

	return ret != 0
}

func GetSystemTimes(lpIdleTime, lpKernelTime, lpUserTime *FILETIME) bool {
	ret, _, _ := procGetSystemTimes.Call(
		uintptr(unsafe.Pointer(lpIdleTime)),
		uintptr(unsafe.Pointer(lpKernelTime)),
		uintptr(unsafe.Pointer(lpUserTime)))

	return ret != 0
}

func GetProcessTimes(hProcess HANDLE, lpCreationTime, lpExitTime, lpKernelTime, lpUserTime *FILETIME) bool {
	ret, _, _ := procGetProcessTimes.Call(
		uintptr(hProcess),
		uintptr(unsafe.Pointer(lpCreationTime)),
		uintptr(unsafe.Pointer(lpExitTime)),
		uintptr(unsafe.Pointer(lpKernelTime)),
		uintptr(unsafe.Pointer(lpUserTime)))

	return ret != 0
}

func GetConsoleScreenBufferInfo(hConsoleOutput HANDLE) *CONSOLE_SCREEN_BUFFER_INFO {
	var csbi CONSOLE_SCREEN_BUFFER_INFO
	ret, _, _ := procGetConsoleScreenBufferInfo.Call(
		uintptr(hConsoleOutput),
		uintptr(unsafe.Pointer(&csbi)))
	if ret == 0 {
		return nil
	}
	return &csbi
}

func SetConsoleTextAttribute(hConsoleOutput HANDLE, wAttributes uint16) bool {
	ret, _, _ := procSetConsoleTextAttribute.Call(
		uintptr(hConsoleOutput),
		uintptr(wAttributes))
	return ret != 0
}

func GetDiskFreeSpaceEx(dirName string) (r bool,
	freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes uint64) {
	ret, _, _ := procGetDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(dirName))),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalNumberOfBytes)),
		uintptr(unsafe.Pointer(&totalNumberOfFreeBytes)))
	return ret != 0,
		freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes
}

func GetSystemTime() *SYSTEMTIME {
	var time SYSTEMTIME
	procGetSystemTime.Call(
		uintptr(unsafe.Pointer(&time)))
	return &time
}

func SetSystemTime(time *SYSTEMTIME) bool {
	ret, _, _ := procSetSystemTime.Call(
		uintptr(unsafe.Pointer(time)))
	return ret != 0
}

func VirtualAllocEx(hProcess HANDLE, lpAddress, dwSize uintptr, flAllocationType, flProtect uint32) uintptr {
	ret, _, _ := procVirtualAllocEx.Call(
		uintptr(hProcess),
		lpAddress,
		dwSize,
		uintptr(flAllocationType),
		uintptr(flProtect),
	)

	return ret
}

func VirtualFreeEx(hProcess HANDLE, lpAddress, dwSize uintptr, dwFreeType uint32) bool {
	ret, _, _ := procVirtualFreeEx.Call(
		uintptr(hProcess),
		lpAddress,
		dwSize,
		uintptr(dwFreeType),
	)

	return ret != 0
}

func WriteProcessMemory(hProcess HANDLE, lpBaseAddress, lpBuffer, nSize uintptr) (int, bool) {
	var nBytesWritten int
	ret, _, _ := procWriteProcessMemory.Call(
		uintptr(hProcess),
		lpBaseAddress,
		lpBuffer,
		nSize,
		uintptr(unsafe.Pointer(&nBytesWritten)),
	)

	return nBytesWritten, ret != 0
}

func ReadProcessMemory(hProcess HANDLE, lpBaseAddress, nSize uintptr) (lpBuffer []uint16, lpNumberOfBytesRead int, ok bool) {

	var nBytesRead int
	buf := make([]uint16, nSize)
	ret, _, _ := procReadProcessMemory.Call(
		uintptr(hProcess),
		lpBaseAddress,
		uintptr(unsafe.Pointer(&buf[0])),
		nSize,
		uintptr(unsafe.Pointer(&nBytesRead)),
	)

	return buf, nBytesRead, ret != 0
}

func QueryPerformanceCounter() uint64 {
	result := uint64(0)
	procQueryPerformanceCounter.Call(
		uintptr(unsafe.Pointer(&result)),
	)

	return result
}

func QueryPerformanceFrequency() uint64 {
	result := uint64(0)
	procQueryPerformanceFrequency.Call(
		uintptr(unsafe.Pointer(&result)),
	)

	return result
}
