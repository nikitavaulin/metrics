package repository

type Storage interface {
	Add(name string, metric any) error
	Get(name string) (any, error)
	GetOrNil(name string) (any, error)
	Update(name string, updated any) error
}
