package pingservice

import (
	"context"

	postgresstorage "github.com/nikitavaulin/metrics/internal/repository/postgres"
)

type PingService struct {
	storage *postgresstorage.Storage
}

func New(storage *postgresstorage.Storage) *PingService {
	return &PingService{storage: storage}
}

func (s *PingService) PingDB(ctx context.Context) error {
	return s.storage.PingDB(ctx)
}
