package products

// Middleware for Logic and HTTP

import (
	"log"
	"net/http"
	"strconv"

	"github.com/wz-Tan/go-ecommerce-api/internal/jsonResponse"
)

type handler struct {
	service Service
}

// Init Handler
func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

// Getting Products
func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())

	// Early Exit
	if err != nil {
		log.Println("Error listing products", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Self-Defined Function
	jsonResponse.Write(w, http.StatusOK, products)
}

// Get Product By ID
func (h *handler) FindProductById(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id") // Get the ID from URL {id}
	id, err := strconv.ParseInt(idString, 10, 64)

	product, err := h.service.FindProductById(r.Context(), id)

	// Early Exit
	if err != nil {
		log.Println("Error listing products", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Self-Defined Function
	jsonResponse.Write(w, http.StatusOK, product)
}
