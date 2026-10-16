package pingservice

import (
	"context"
	"errors"
)

type Pinger interface {
	PingDB(context.Context) error
}

type PingService struct {
	storage Pinger
}

func New(storage Pinger) *PingService {
	return &PingService{storage: storage}
}

func (s *PingService) PingDB(ctx context.Context) error {
	if s.storage == nil {
		return errors.New("postgres is not configured")
	}
	return s.storage.PingDB(ctx)
}
