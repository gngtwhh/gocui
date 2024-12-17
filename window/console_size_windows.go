//go:build windows || !unix

package window

import (
	"syscall"
	"unsafe"
)

type (
	short     int16
	word      uint16
	ulong     uint32
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
	/*
		typedef struct _CONSOLE_SCREEN_BUFFER_INFO {
		  COORD      dwSize;
		  COORD      dwCursorPosition;
		  WORD       wAttributes;
		  SMALL_RECT srWindow;
		  COORD      dwMaximumWindowSize;
		} CONSOLE_SCREEN_BUFFER_INFO;
	*/
	consoleScreenBufferInfo struct {
		dwSize              coord
		dwCursorPosition    coord
		wAttributes         word
		srWindow            smallRect
		dwMaximumWindowSize coord
	}
	/*
		typedef struct _CONSOLE_SCREEN_BUFFER_INFOEX {
		  ULONG      cbSize;
		  COORD      dwSize;
		  COORD      dwCursorPosition;
		  WORD       wAttributes;
		  SMALL_RECT srWindow;
		  COORD      dwMaximumWindowSize;
		  WORD       wPopupAttributes;
		  BOOL       bFullscreenSupported;
		  COLORREF   ColorTable[16];
		} CONSOLE_SCREEN_BUFFER_INFOEX, *PCONSOLE_SCREEN_BUFFER_INFOEX;
	*/
	consoleScreenBufferInfoEx struct {
		cbSize               uint32
		dwSize               coord
		dwCursorPosition     coord
		wAttributes          word
		srWindow             smallRect
		dwMaximumWindowSize  coord
		wPopupAttributes     word
		bFullscreenSupported bool
		ColorTable           [16]uint32
	}
)

func (c coord) uintptr() uintptr {
	// little endian, put x first
	return uintptr(c.X) | (uintptr(c.Y) << 16)
}

var kernel32DLL = syscall.NewLazyDLL("kernel32.dll")
var getConsoleScreenBufferInfoProc = kernel32DLL.NewProc("GetConsoleScreenBufferInfo")

//var setConsoleScreenBufferInfoProc = kernel32DLL.NewProc("SetConsoleScreenBufferInfo")
//var setConsoleWindowInfoProc = kernel32DLL.NewProc("SetConsoleWindowInfo")

var handle *syscall.Handle

func getConsoleScreenBufferInfo() (*consoleScreenBufferInfo, error) {
	if handle == nil {
		h, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
		if err != nil {
			return nil, err
		}
		handle = &h
	}

	var info consoleScreenBufferInfo
	if err := getError(getConsoleScreenBufferInfoProc.Call(uintptr(*handle), uintptr(unsafe.Pointer(&info)))); err != nil {
		return nil, err
	}
	return &info, nil
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

func GetConsoleSize() (weight, height int) {
	info, err := getConsoleScreenBufferInfo()
	if err != nil {
		return 0, 0
	}
	return int(info.srWindow.Right - info.srWindow.Left + 1), int(info.srWindow.Bottom - info.srWindow.Top + 1)
}
