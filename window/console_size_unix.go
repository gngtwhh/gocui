//go:build unix || !windows

package window

import (
	"syscall"
	"unsafe"
)

func GetConsoleSize() (weight, height int) {
	var sz struct {
		rows   uint16
		cols   uint16
		xpixel uint16
		ypixel uint16
	}
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdout), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&sz)))
	return int(sz.cols), int(sz.rows)
}
