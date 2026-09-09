package orders

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	repo "github.com/wz-Tan/go-ecommerce-api/internal/adapters/postgresql/sqlc"
)

type Service interface {
	CreateOrder(ctx context.Context, op OrderPayload) (repo.Order, error)
}

type svc struct {
	repo *repo.Queries
	db   *pgx.Conn
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) CreateOrder(ctx context.Context, op OrderPayload) (repo.Order, error) {
	// Verify Customer
	if op.CustomerID == 0 {
		return repo.Order{}, fmt.Errorf("Invalid Customer ID")
	}

	// Verify Product Exists
	product, err := s.repo.FindProductById(ctx, op.ProductID)
	if err != nil {
		return repo.Order{}, fmt.Errorf("Invalid Product ID")
	}

	// Verify Enough Stock
	if op.Quantity > product.Quantity || op.Quantity <= 0 {
		return repo.Order{}, fmt.Errorf("Insufficient Product Quantity")
	}

	// Connect to DB for Commit and Rollback (tx -> transaction)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, fmt.Errorf("Can't start database instance!")
	}

	// Undo In The End
	defer tx.Rollback(ctx)

	// Repo Uses Current DB Instance
	qtx := s.repo.WithTx(tx)

	// Create Order Record, Then Return Order
	order, err := qtx.CreateOrder(ctx, op.CustomerID)
	if err != nil {
		return repo.Order{}, fmt.Errorf("Error Creating Order")
	}

	// Create Order Item Record
	createOrderItemParams := repo.CreateOrderItemParams{
		OrderID:    order.ID,
		ProductID:  op.ProductID,
		Quantity:   op.Quantity,
		PriceCents: op.PriceCents,
	}

	// Create Join Table
	err = qtx.CreateOrderItem(ctx, createOrderItemParams)
	if err != nil {
		return repo.Order{}, fmt.Errorf("Error Creating Order Item")
	}

	// Update Stock Number
	err = qtx.UpdateProductQuantity(ctx, repo.UpdateProductQuantityParams{
		ID:       op.ProductID,
		Quantity: product.Quantity - op.Quantity,
	})

	if err != nil {
		return repo.Order{}, fmt.Errorf("Error Updating Product Quantity")
	}

	// Commit
	if err := tx.Commit(ctx); err != nil {
		return repo.Order{}, fmt.Errorf("Error Committing DB!")
	}

	// No Issues
	return order, nil
}
