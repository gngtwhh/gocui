//go:build windows || !unix

package window

import (
	"fmt"
	"syscall"
	"unsafe"
)

type (
	short     int16
	word      uint16
	smallRect struct {
		Left   short
		Top    short
		Right  short
		Bottom short
	}
	coord struct {
		X short
		Y short
	}
	lpConsoleScreenBufferInfo struct {
		dwSize              coord
		dwCursorPosition    coord
		wAttributes         word
		srWindow            smallRect
		dwMaximumWindowSize coord
	}
)

var kernel32DLL = syscall.NewLazyDLL("kernel32.dll")
var getConsoleScreenBufferInfoProc = kernel32DLL.NewProc("GetConsoleScreenBufferInfo")

func GetConsoleSize() (weight, height int) {
	handle, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	if err != nil {
		panic(fmt.Errorf("could not get std I/O handle"))
	}

	var info lpConsoleScreenBufferInfo
	if err = getError(getConsoleScreenBufferInfoProc.Call(uintptr(handle), uintptr(unsafe.Pointer(&info)))); err != nil {
		return 0, 0
	}

	return int(info.srWindow.Right - info.srWindow.Left + 1), int(info.srWindow.Bottom - info.srWindow.Top + 1)
}

func getError(r1, r2 uintptr, lastErr error) error {
	// If the function fails, the return value is zero.
	if r1 == 0 {
		if lastErr != nil {
			return lastErr
		}
		return syscall.EINVAL
	}
	return nil
}
