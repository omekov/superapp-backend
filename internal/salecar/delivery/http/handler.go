package http

import (
	"net/http"

	"github.com/omekov/superapp-backend/internal/salecar/car/service"
)

// Handler ...
type Handler struct {
	service *service.Service
}

// NewHandler ...
func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Init ...
func (h *Handler) Init() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return mux
}
