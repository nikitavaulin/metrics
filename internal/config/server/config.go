package serverconfig

type HTTPServerConfig struct {
	Address string
}

func New() (*HTTPServerConfig, error) {
	return &HTTPServerConfig{
		Address: ":8080",
	}, nil
}
