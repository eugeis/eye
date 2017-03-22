package integ

type Cache interface {
	Get(key string, builder func() interface{}) (value interface{}, ok bool)
	GetOrBuild(key string, builder func() (interface{}, error)) (value interface{}, err error)
	Put(key string, value interface{})
	Clear()
}
