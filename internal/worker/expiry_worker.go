package worker

import (
	"context"
	"log"
	"time"

	"github.com/realtime-auction/internal/repository"
)

type ExpiryWorker struct {
	auctionRepo *repository.AuctionRepository
	interval    time.Duration
}

func NewExpiryWorker(auctionRepo *repository.AuctionRepository, interval time.Duration) *ExpiryWorker {
	return &ExpiryWorker{
		auctionRepo: auctionRepo,
		interval:    interval,
	}
}

func (w *ExpiryWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.processExpiredAuctions(ctx)
		}
	}
}

func (w *ExpiryWorker) processExpiredAuctions(ctx context.Context) {
	auctions, err := w.auctionRepo.GetExpiredActiveAuctions(ctx)
	if err != nil {
		log.Printf("expiry worker: failed to get expired auctions: %v", err)
		return
	}

	for _, a := range auctions {
		if err := w.auctionRepo.UpdateStatus(ctx, a.ID, "ended"); err != nil {
			log.Printf("expiry worker: failed to end auction %d: %v", a.ID, err)
			continue
		}
		log.Printf("expiry worker: ended auction %d (%s)", a.ID, a.Title)
	}
}
