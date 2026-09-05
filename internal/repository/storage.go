package repository

type Storage interface {
	Add(name string, metric any) error
	Get(name string) (any, error)
	GetList() (map[string]any, error)
	GetOrNil(name string) (any, error)
	Update(name string, updated any) error
}

type FileStorage interface {
	SaveToJSON(filename string, data any) error
	LoadFromJSON(filename string, dest any) error
}
