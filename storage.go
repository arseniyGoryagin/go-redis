package main

import "sync"

var storageSet = map[string]string{}
var storageSetMutex = sync.RWMutex{}

func Set(key string, val string) {
	storageSetMutex.Lock()
	defer storageSetMutex.Unlock()
	storageSet[key] = val
}

func Get(key string) (bool, string) {
	storageSetMutex.RLock()
	defer storageSetMutex.RUnlock()
	val, ok := storageSet[key]
	return ok, val
}

var storageHset = map[string]map[string]string{}
var storageHsetMutex = sync.RWMutex{}

func Hset(hash string, key string, val string) {
	storageHsetMutex.Lock()
	defer storageHsetMutex.Unlock()
	_, ok := storageHset[hash]
	if !ok {
		storageHset[hash] = map[string]string{}
	}
	storageHset[hash][key] = val
}

func Hget(hash string, key string) (bool, string) {
	storageHsetMutex.Lock()
	defer storageHsetMutex.Unlock()
	_, ok := storageHset[hash]
	if !ok {
		return false, ""
	}
	return true, storageHset[hash][key]
}
