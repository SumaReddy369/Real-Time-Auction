package models

import "time"

type Bid struct {
	ID        int64     `json:"id"`
	AuctionID int64     `json:"auction_id"`
	BidderID  string    `json:"bidder_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type PlaceBidRequest struct {
	BidderID string  `json:"bidder_id"`
	Amount   float64 `json:"amount"`
}
