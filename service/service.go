package service

type Service interface {
	Name() string
	Kind() string

	Init() error
	Close()

	Ping() error

	NewСheck(query string, expr string) (Check, error)
}

type Check interface {
	Check() (bool, error)
}

type Factory interface {
	Find(name string) Service
	Close()
}