package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/wz-Tan/go-ecommerce-api/internal/env"
)

func main() {
	// Structured Logging (Logs as Dicts) and In Output Channel
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger) // Replaces the OG Log

	ctx := context.Background()

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=postgres dbname=go_ecommerce sslmode=disable")},
	}

	// Connect to DB
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx) // Run At the End (Cleanup)

	logger.Info("connected to database", "dsn", cfg.db.dsn)

	api := application{
		config: cfg,
		db:     conn,
	}

	// Handler (Both under API as an umbrella call)
	h := api.mount()

	// Run and Serve
	if err := api.run(h); err != nil { // Assign; Check
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
