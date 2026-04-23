package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dulatAsisADV2/order-service/internal/domain"
	"dulatAsisADV2/order-service/internal/usecase"
)

type OrderHandler struct {
	orderUseCase *usecase.OrderUseCase
}

func NewOrderHandler(orderUseCase *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
	}
}

type CreateOrderRequest struct {
	UserID      string  `json:"user_id"`
	ProductID   string  `json:"product_id"`
	Quantity    int32   `json:"quantity"`
	Price       float64 `json:"price"`
	TotalAmount float64 `json:"total_amount"`
}

type OrderResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Status      string  `json:"status"`
	TotalAmount float64 `json:"total_amount"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create order items
	items := []*domain.OrderItem{
		{
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			Price:     req.Price,
		},
	}

	order, err := h.orderUseCase.CreateOrder(c.Request.Context(), req.UserID, items, req.TotalAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
	})
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	order, err := h.orderUseCase.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
	})
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.orderUseCase.UpdateOrderStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order status updated"})
}

func NewRouter(orderUseCase *usecase.OrderUseCase) *gin.Engine {
	router := gin.Default()

	handler := NewOrderHandler(orderUseCase)

	api := router.Group("/api/orders")
	{
		api.POST("/", handler.CreateOrder)
		api.GET("/:id", handler.GetOrder)
		api.PUT("/:id/status", handler.UpdateOrderStatus)
	}

	return router
}
