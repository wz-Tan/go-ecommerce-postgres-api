package orders

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wz-Tan/go-ecommerce-api/internal/jsonResponse"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Decode Request to get payload
	var p OrderPayload

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		jsonResponse.Write(w, http.StatusBadRequest, "Invalid JSON Request Body")
		return
	}

	defer r.Body.Close()

	// Create Order and Order Item, Then Get Order
	order, err := h.service.CreateOrder(r.Context(), p)
	if err != nil {
		jsonResponse.Write(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Order Created: %+v", order)
}
