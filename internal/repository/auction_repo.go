package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/realtime-auction/internal/models"
	"time"
)

type AuctionRepository struct {
	pool *pgxpool.Pool
}

func NewAuctionRepository(pool *pgxpool.Pool) *AuctionRepository {
	return &AuctionRepository{pool: pool}
}

func (r *AuctionRepository) Create(ctx context.Context, req *models.CreateAuctionRequest) (*models.Auction, error) {
	endTime := time.Now().Add(time.Duration(req.DurationMin) * time.Minute)
	var a models.Auction
	err := r.pool.QueryRow(ctx, `
		INSERT INTO auctions (title, description, start_price, current_bid, end_time, status, created_at, updated_at)
		VALUES ($1, $2, $3, $3, $4, 'active', NOW(), NOW())
		RETURNING id, title, description, start_price, current_bid, end_time, status, created_at, updated_at
	`, req.Title, req.Description, req.StartPrice, endTime).Scan(
		&a.ID, &a.Title, &a.Description, &a.StartPrice, &a.CurrentBid, &a.EndTime, &a.Status, &a.CreatedAt, &a.UpdatedAt,
	)
	return &a, err
}

func (r *AuctionRepository) GetByID(ctx context.Context, id int64) (*models.Auction, error) {
	var a models.Auction
	err := r.pool.QueryRow(ctx, `
		SELECT id, title, description, start_price, current_bid, end_time, status, created_at, updated_at
		FROM auctions WHERE id = $1
	`, id).Scan(&a.ID, &a.Title, &a.Description, &a.StartPrice, &a.CurrentBid, &a.EndTime, &a.Status, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AuctionRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE auctions SET status = $1, updated_at = NOW() WHERE id = $2`, status, id)
	return err
}

func (r *AuctionRepository) UpdateCurrentBid(ctx context.Context, id int64, amount float64) error {
	_, err := r.pool.Exec(ctx, `UPDATE auctions SET current_bid = $1, updated_at = NOW() WHERE id = $2`, amount, id)
	return err
}

func (r *AuctionRepository) GetExpiredActiveAuctions(ctx context.Context) ([]*models.Auction, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, description, start_price, current_bid, end_time, status, created_at, updated_at
		FROM auctions WHERE status = 'active' AND end_time <= NOW()
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var auctions []*models.Auction
	for rows.Next() {
		var a models.Auction
		if err := rows.Scan(&a.ID, &a.Title, &a.Description, &a.StartPrice, &a.CurrentBid, &a.EndTime, &a.Status, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		auctions = append(auctions, &a)
	}
	return auctions, rows.Err()
}
