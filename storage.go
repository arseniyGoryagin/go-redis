package main

import "sync"

var storageSet = map[string]string{}
var storageSetMutex = sync.RWMutex{}

func Set(key string, val string) {
	storageSetMutex.Lock()
	defer storageSetMutex.Unlock()
	storageSet[key] = val
}

func Get(key string) (string, bool) {
	storageSetMutex.RLock()
	defer storageSetMutex.RUnlock()
	val, ok := storageSet[key]
	return val, ok
}
