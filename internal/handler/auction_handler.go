package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/realtime-auction/internal/models"
	"github.com/realtime-auction/internal/repository"
	"github.com/realtime-auction/internal/redis"
)

type AuctionHandler struct {
	auctionRepo *repository.AuctionRepository
	bidRepo     *repository.BidRepository
	pubsub      *redis.PubSub
}

func NewAuctionHandler(auctionRepo *repository.AuctionRepository, bidRepo *repository.BidRepository, pubsub *redis.PubSub) *AuctionHandler {
	return &AuctionHandler{
		auctionRepo: auctionRepo,
		bidRepo:     bidRepo,
		pubsub:      pubsub,
	}
}

func (h *AuctionHandler) CreateAuction(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAuctionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Title == "" || req.StartPrice <= 0 || req.DurationMin <= 0 {
		respondError(w, http.StatusBadRequest, "title, start_price, and duration_min are required")
		return
	}

	auction, err := h.auctionRepo.Create(r.Context(), &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create auction")
		return
	}
	respondJSON(w, http.StatusCreated, auction)
}

func (h *AuctionHandler) GetAuction(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid auction id")
		return
	}

	auction, err := h.auctionRepo.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "auction not found")
		return
	}
	respondJSON(w, http.StatusOK, auction)
}

func (h *AuctionHandler) PlaceBid(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid auction id")
		return
	}

	var req models.PlaceBidRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.BidderID == "" || req.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "bidder_id and amount are required")
		return
	}

	auction, err := h.auctionRepo.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "auction not found")
		return
	}
	if auction.Status != "active" {
		respondError(w, http.StatusBadRequest, "auction is not active")
		return
	}
	if req.Amount <= auction.CurrentBid {
		respondError(w, http.StatusBadRequest, "bid must be higher than current bid")
		return
	}

	bid, err := h.bidRepo.Create(r.Context(), id, req.BidderID, req.Amount)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to place bid")
		return
	}

	if err := h.auctionRepo.UpdateCurrentBid(r.Context(), id, req.Amount); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update auction")
		return
	}

	// Broadcast to WebSocket subscribers
	_ = h.pubsub.PublishBid(r.Context(), id, &redis.BidMessage{
		AuctionID: id,
		BidderID:  req.BidderID,
		Amount:    req.Amount,
	})

	respondJSON(w, http.StatusCreated, bid)
}

func (h *AuctionHandler) ListBids(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid auction id")
		return
	}

	_, err = h.auctionRepo.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "auction not found")
		return
	}

	bids, err := h.bidRepo.ListByAuctionID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list bids")
		return
	}
	if bids == nil {
		bids = []*models.Bid{}
	}
	respondJSON(w, http.StatusOK, bids)
}
