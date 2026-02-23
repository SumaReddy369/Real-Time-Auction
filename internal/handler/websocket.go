package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/realtime-auction/internal/redis"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	pubsub *redis.PubSub
	clients map[int64]map[*websocket.Conn]struct{}
	mu      sync.RWMutex
}

func NewWebSocketHandler(pubsub *redis.PubSub) *WebSocketHandler {
	return &WebSocketHandler{
		pubsub:  pubsub,
		clients: make(map[int64]map[*websocket.Conn]struct{}),
	}
}

func (h *WebSocketHandler) Handle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid auction id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	h.register(id, conn)
	defer h.unregister(id, conn)

	sub, err := h.pubsub.Subscribe(r.Context(), id)
	if err != nil {
		log.Printf("subscribe error: %v", err)
		return
	}
	defer sub.Close()

	ch := sub.Channel()
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			var bidMsg redis.BidMessage
			if err := json.Unmarshal([]byte(msg.Payload), &bidMsg); err != nil {
				continue
			}
			if err := conn.WriteJSON(bidMsg); err != nil {
				return
			}
		}
	}
}

func (h *WebSocketHandler) register(auctionID int64, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[auctionID] == nil {
		h.clients[auctionID] = make(map[*websocket.Conn]struct{})
	}
	h.clients[auctionID][conn] = struct{}{}
}

func (h *WebSocketHandler) unregister(auctionID int64, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[auctionID] != nil {
		delete(h.clients[auctionID], conn)
		if len(h.clients[auctionID]) == 0 {
			delete(h.clients, auctionID)
		}
	}
}
