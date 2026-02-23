-- Auctions table
CREATE TABLE IF NOT EXISTS auctions (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_price DECIMAL(12, 2) NOT NULL,
    current_bid DECIMAL(12, 2) NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Bids table
CREATE TABLE IF NOT EXISTS bids (
    id BIGSERIAL PRIMARY KEY,
    auction_id BIGINT NOT NULL REFERENCES auctions(id) ON DELETE CASCADE,
    bidder_id VARCHAR(100) NOT NULL,
    amount DECIMAL(12, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_bids_auction_id ON bids(auction_id);
CREATE INDEX idx_auctions_status_end_time ON auctions(status, end_time);
