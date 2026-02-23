package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"

	"github.com/realtime-auction/config"
	"github.com/realtime-auction/internal/db"
	"github.com/realtime-auction/internal/handler"
	"github.com/realtime-auction/internal/redis"
	"github.com/realtime-auction/internal/repository"
	"github.com/realtime-auction/internal/worker"
)

func main() {
	cfg := config.Load()

	// PostgreSQL
	pool, err := pgxpool.New(context.Background(), cfg.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("connected to PostgreSQL")

	if err := db.Migrate(context.Background(), pool); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("migrations completed")

	// Redis
	rdb := goredis.NewClient(&goredis.Options{
		Addr: cfg.RedisAddr,
		DB:   cfg.RedisDB,
	})
	defer rdb.Close()

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	log.Println("connected to Redis")

	// Repositories
	auctionRepo := repository.NewAuctionRepository(pool)
	bidRepo := repository.NewBidRepository(pool)

	// Redis PubSub
	pubsub := redis.NewPubSub(rdb)

	// Handlers
	auctionHandler := handler.NewAuctionHandler(auctionRepo, bidRepo, pubsub)
	wsHandler := handler.NewWebSocketHandler(pubsub)

	// Router
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/auctions", auctionHandler.CreateAuction).Methods("POST")
	r.HandleFunc("/api/v1/auctions/{id}", auctionHandler.GetAuction).Methods("GET")
	r.HandleFunc("/api/v1/auctions/{id}/bids", auctionHandler.PlaceBid).Methods("POST")
	r.HandleFunc("/api/v1/auctions/{id}/bids", auctionHandler.ListBids).Methods("GET")
	r.HandleFunc("/ws/auctions/{id}", wsHandler.Handle)

	// CORS
	r.Use(corsMiddleware)

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Background worker
	expiryWorker := worker.NewExpiryWorker(auctionRepo, 10*time.Second)
	go expiryWorker.Run(context.Background())

	// Server
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("server listening on :%s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	fmt.Println("server stopped")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
