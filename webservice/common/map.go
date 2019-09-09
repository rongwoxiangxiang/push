package common

import "sync"

type Map struct {
	sync.RWMutex
	mp map[string]interface{}
}

func (this *Map) Set(key string, val interface{}) {
	this.Lock()
	defer this.Unlock()
	this.mp[key] = val
}

func (this *Map) Get(key string) interface{} {
	this.RLock()
	defer this.RUnlock()

	if val, ok := this.mp[key]; ok {
		return val
	}
	return nil
}

func (this *Map) Del(key string) {
	this.Lock()
	defer this.Unlock()
	delete(this.mp, key)
}

func (this *Map) Count() int {
	return len(this.mp)
}

func (this *Map) List() map[string]interface{} {
	return this.mp
}
