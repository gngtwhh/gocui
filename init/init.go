package init

import "runtime"

const (
	Windows = iota
	Linux
	Unix
	Mac
)

var SysType int

func init() {
	switch runtime.GOOS {
	case "windows":
		SysType = Windows
	case "Linux":
		SysType = Linux
	}
}
