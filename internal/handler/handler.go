package handler

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"io"
	"justpay/internal/config"
	"justpay/internal/service"
	"net/http"
)

var upgrader websocket.Upgrader

type Handler struct {
	cfg     *config.Config
	service service.Service
}

func New(cfg *config.Config) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

func (h *Handler) Init() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /webhooks/payments/orders", h.OrdersPaymentsWebHook)
	mux.HandleFunc("GET /orders/{order_id}/events", h.OrderEventsHandler)
	mux.HandleFunc("GET /orders", h.OrdersHandler)

	return mux
}

type OrdersPaymentsRequest struct {
}

func (h *Handler) OrdersPaymentsWebHook(w http.ResponseWriter, r *http.Request) {
	var req OrdersPaymentsRequest

	rawReq, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO
		return
	}

	if err = json.Unmarshal(rawReq, &req); err != nil {
		// TODO
		return
	}

}

func (h *Handler) OrderEventsHandler(w http.ResponseWriter, r *http.Request) {
	upgrade, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	orderID := r.URL.Query().Get("order_id")
	if orderID == "" {
		// TODO
		return
	}

	h.service.Order.Subscribe()

	for {
		upgrade.WriteJSON()
	}
}

func (h *Handler) OrdersHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
