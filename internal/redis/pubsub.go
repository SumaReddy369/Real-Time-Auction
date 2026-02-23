package redis

import (
	"context"
	"encoding/json"
	"strconv"

	goredis "github.com/redis/go-redis/v9"
)

const (
	ChannelPrefix = "auction:"
)

type BidMessage struct {
	AuctionID int64   `json:"auction_id"`
	BidderID  string  `json:"bidder_id"`
	Amount    float64 `json:"amount"`
}

type PubSub struct {
	client *goredis.Client
}

func NewPubSub(client *goredis.Client) *PubSub {
	return &PubSub{client: client}
}

func (p *PubSub) PublishBid(ctx context.Context, auctionID int64, msg *BidMessage) error {
	channel := ChannelPrefix + formatAuctionID(auctionID)
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.client.Publish(ctx, channel, data).Err()
}

func (p *PubSub) Subscribe(ctx context.Context, auctionID int64) (*goredis.PubSub, error) {
	channel := ChannelPrefix + formatAuctionID(auctionID)
	ps := p.client.Subscribe(ctx, channel)
	_, err := ps.Receive(ctx)
	return ps, err
}

func formatAuctionID(id int64) string {
	return strconv.FormatInt(id, 10)
}
