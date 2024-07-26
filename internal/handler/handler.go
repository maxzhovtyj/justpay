package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"justpay/internal/config"
	"justpay/internal/domain/order"
	"justpay/internal/service"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	cfg     *config.Config
	service *service.Service
}

func New(cfg *config.Config, s *service.Service) *Handler {
	return &Handler{
		cfg:     cfg,
		service: s,
	}
}

func (h *Handler) Init() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /webhooks/payments/orders", h.OrdersPaymentsWebHook)
	mux.HandleFunc("GET /orders/{order_id}/events", h.OrderEventsHandler)
	mux.HandleFunc("GET /orders", h.OrdersHandler)

	return mux
}

func httpResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

type OrdersPaymentsRequest struct {
	EventId     string    `json:"event_id"`
	OrderId     string    `json:"order_id"`
	UserId      string    `json:"user_id"`
	OrderStatus string    `json:"order_status"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (h *Handler) OrdersPaymentsWebHook(w http.ResponseWriter, r *http.Request) {
	var req OrdersPaymentsRequest

	rawReq, err := io.ReadAll(r.Body)
	if err != nil {
		httpResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = json.Unmarshal(rawReq, &req); err != nil {
		httpResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	event, err := order.NewEvent(req.EventId, req.OrderId, req.UserId, req.OrderStatus, req.CreatedAt, req.UpdatedAt)
	if err != nil {
		httpResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Order.NewEvent(event)
	if err != nil {
		if errors.Is(err, order.ErrEventAlreadyExists) {
			httpResponse(w, http.StatusConflict, err.Error())
			return
		}

		if errors.Is(err, order.ErrFinalStatusReceived) {
			httpResponse(w, http.StatusGone, err.Error())
			return
		}

		httpResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpResponse(w, http.StatusOK, "successful")
}

func (h *Handler) OrderEventsHandler(w http.ResponseWriter, r *http.Request) {
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	// SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	rawOrderID := r.PathValue("order_id")
	if rawOrderID == "" {
		httpResponse(w, http.StatusBadRequest, "empty order_id")
		return
	}

	orderID, err := uuid.Parse(rawOrderID)
	if err != nil {
		httpResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	sub, err := h.service.Order.Subscribe(orderID)
	if err != nil {
		httpResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer h.service.Order.Unsubscribe(orderID, sub.ID)

	rc := http.NewResponseController(w)

	log.Printf("%s start listening to events", sub.ID.String())
	for e := range sub.ReadEvents() {
		rawMessage, err := json.Marshal(e)
		if err != nil {
			log.Printf("error marshalling message: %v", err)
			continue
		}

		_, err = fmt.Fprintf(w, "data: %s\n\n", rawMessage)
		if err != nil {
			log.Printf("error writing message: %v", err)
			continue
		}

		err = rc.Flush()
		if err != nil {
			log.Printf("error flushing message: %v", err)
			continue
		}
	}
}

func (h *Handler) OrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	orders, err := h.service.Order.GetOrders()
	if err != nil {
		httpResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	rawOrders, err := json.Marshal(orders)
	if err != nil {
		httpResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = w.Write(rawOrders)
	if err != nil {
		httpResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
