package w32

import (
	"syscall"

	"golang.org/x/sys/windows"
)

var (
	modmsimg32 = windows.NewLazySystemDLL("msimg32.dll")

	procAlphaBlend = modmsimg32.NewProc("AlphaBlend")
)

func AlphaBlend(dcdest HDC, xoriginDest int32, yoriginDest int32, wDest int32, hDest int32, dcsrc HDC, xoriginSrc int32, yoriginSrc int32, wsrc int32, hsrc int32, ftn uintptr) (err error) {
	r1, _, e1 := syscall.Syscall12(procAlphaBlend.Addr(), 11, uintptr(dcdest), uintptr(xoriginDest), uintptr(yoriginDest), uintptr(wDest), uintptr(hDest), uintptr(dcsrc), uintptr(xoriginSrc), uintptr(yoriginSrc), uintptr(wsrc), uintptr(hsrc), uintptr(ftn), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
