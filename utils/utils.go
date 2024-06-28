package utils

import "sync"

var ConsoleMutex sync.Mutex

func init() {
	ConsoleMutex = sync.Mutex{}
}
