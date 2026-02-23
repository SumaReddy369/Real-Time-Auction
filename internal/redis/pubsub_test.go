package redis

import (
	"testing"
)

func TestBidMessage(t *testing.T) {
	msg := &BidMessage{
		AuctionID: 1,
		BidderID:  "user1",
		Amount:    100.50,
	}
	if msg.AuctionID != 1 || msg.BidderID != "user1" || msg.Amount != 100.50 {
		t.Error("BidMessage fields incorrect")
	}
}

// Note: PublishBid and Subscribe require a real Redis connection
// Run integration tests with docker-compose for full coverage
