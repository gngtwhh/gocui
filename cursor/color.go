package cursor

import "fmt"

const (
	RESET   = "\033[0m"
	RED     = "\033[31m"
	GREEN   = "\033[32m"
	YELLOW  = "\033[33m"
	BLUE    = "\033[34m"
	MAGENTA = "\033[35m"
	CYAN    = "\033[36m"
	WHITE   = "\033[37m"
)

func SetColor(color string) {
	fmt.Printf("%s", color)
}

func ResetColor() {
	fmt.Printf("%s", RESET)
}
