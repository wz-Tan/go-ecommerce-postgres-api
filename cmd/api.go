package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	repo "github.com/wz-Tan/go-ecommerce-api/internal/adapters/postgresql/sqlc"
	"github.com/wz-Tan/go-ecommerce-api/internal/orders"
	"github.com/wz-Tan/go-ecommerce-api/internal/products"
)

// functions -> run, mount, shutdown instance

// mounter (* means it modifies itself and not a new object)
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID) // for rate limiting
	r.Use(middleware.RealIP)    // for rate limiting, analytics and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // crash recoverer

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("No problem!"))
	})

	querier := repo.New(app.db)

	// Product Handler -> Middleware between Logic and HTTP
	productService := products.NewService(querier)
	productHandler := products.NewHandler(productService)

	// Feed in Response and Request by Default (Response is Returned in a Stream and ends as function ends)
	r.Get("/products", productHandler.ListProducts)

	// Get Product by ID
	r.Get("/products/{id}", productHandler.FindProductById)

	// Order Management
	orderService := orders.NewService(querier, app.db)
	orderHandler := orders.NewHandler(orderService)

	// Post an Order
	r.Post("/createOrder", orderHandler.CreateOrder)

	return r
}

// run (function belongs to application, takes in handler and returns error)
func (app *application) run(h http.Handler) error {
	// Server Consists of App Config and Handler
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute * 1,
	}

	log.Printf("Server has started at address %s", srv.Addr)

	return srv.ListenAndServe()
}

// server
type application struct {
	config config
	db     *pgx.Conn
}

// server config
type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}
