package pingservice

import (
	"context"
	"errors"

	postgresstorage "github.com/nikitavaulin/metrics/internal/repository/postgres"
)

type PingService struct {
	storage *postgresstorage.Storage
}

func New(storage *postgresstorage.Storage) *PingService {
	return &PingService{storage: storage}
}

func (s *PingService) PingDB(ctx context.Context) error {
	if s.storage == nil {
		return errors.New("postgres is not configured")
	}
	return s.storage.PingDB(ctx)
}
