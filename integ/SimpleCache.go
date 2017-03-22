package integ

import (
	"sync"
)


type SimpleCache struct {
	MaxSize int
	data    map[string]interface{}
	lock    sync.Mutex
}

func NewCache() SimpleCache {
	return SimpleCache{MaxSize: 1000, data: make(map[string]interface{}), lock: sync.Mutex{}}
}

func (o SimpleCache) Clear() {
	o.lock.Lock()
	o.data = make(map[string]interface{})
	o.lock.Unlock()
	return
}

func (o SimpleCache) Get(key string, builder func() interface{}) (value interface{}, ok bool) {
	o.lock.Lock()
	value, ok = o.data[key]
	o.lock.Unlock()
	return
}

func (o SimpleCache) GetOrBuild(key string, builder func() (interface{}, error)) (value interface{}, err error) {
	o.lock.Lock()
	value, ok := o.data[key]
	if !ok {
		value, err = builder()
		if err == nil {
			o.put(key, value)
		}
	}
	o.lock.Unlock()
	return
}

func (o SimpleCache) Put(key string, value interface{}) {
	o.lock.Lock()
	o.put(key, value)
	o.lock.Unlock()
}

func (o SimpleCache) put(key string, value interface{}) {
	//reset cache
	if len(o.data) >= o.MaxSize {
		o.data = make(map[string]interface{})
	}
	o.data[key] = value
}
