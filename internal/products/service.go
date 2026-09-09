package products

import (
	"context"

	repo "github.com/wz-Tan/go-ecommerce-api/internal/adapters/postgresql/sqlc"
)

// Logic
// context is passed by http with its kill switch
type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	FindProductById(ctx context.Context, id int64) (repo.Product, error)
}

// svc instead of Service cuz name was used
type svc struct {
	// database
	repo repo.Querier // Created when we ran sqlc (contains all queries)
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *svc) FindProductById(ctx context.Context, id int64) (repo.Product, error) {
	return s.repo.FindProductById(ctx, id)
}
