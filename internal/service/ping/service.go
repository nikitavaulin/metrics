package pingservice

import (
	"context"
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
	return s.storage.PingDB(ctx)
}
