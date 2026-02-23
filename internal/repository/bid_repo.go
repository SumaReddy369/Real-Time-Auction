package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/realtime-auction/internal/models"
)

type BidRepository struct {
	pool *pgxpool.Pool
}

func NewBidRepository(pool *pgxpool.Pool) *BidRepository {
	return &BidRepository{pool: pool}
}

func (r *BidRepository) Create(ctx context.Context, auctionID int64, bidderID string, amount float64) (*models.Bid, error) {
	var b models.Bid
	err := r.pool.QueryRow(ctx, `
		INSERT INTO bids (auction_id, bidder_id, amount, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, auction_id, bidder_id, amount, created_at
	`, auctionID, bidderID, amount).Scan(&b.ID, &b.AuctionID, &b.BidderID, &b.Amount, &b.CreatedAt)
	return &b, err
}

func (r *BidRepository) ListByAuctionID(ctx context.Context, auctionID int64) ([]*models.Bid, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, auction_id, bidder_id, amount, created_at
		FROM bids WHERE auction_id = $1 ORDER BY created_at DESC
	`, auctionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []*models.Bid
	for rows.Next() {
		var b models.Bid
		if err := rows.Scan(&b.ID, &b.AuctionID, &b.BidderID, &b.Amount, &b.CreatedAt); err != nil {
			return nil, err
		}
		bids = append(bids, &b)
	}
	return bids, rows.Err()
}
