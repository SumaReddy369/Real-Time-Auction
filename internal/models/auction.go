package models

import "time"

type Auction struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartPrice  float64   `json:"start_price"`
	CurrentBid  float64   `json:"current_bid"`
	EndTime     time.Time `json:"end_time"`
	Status      string    `json:"status"` // active, ended, cancelled
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateAuctionRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	StartPrice  float64 `json:"start_price"`
	DurationMin int     `json:"duration_min"`
}
