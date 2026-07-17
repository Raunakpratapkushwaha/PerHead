package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/api/middleware"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/model"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/service"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(ps service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: ps}
}

func (h *PaymentHandler) RecordPayment(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var req model.RecordPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// USING YOUR CUSTOM MIDDLEWARE
	userID := int64(middleware.GetUserID(c))
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}

	payment, err := h.paymentService.RecordPayment(c.Request.Context(), groupID, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

func (h *PaymentHandler) GetGroupPayments(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	userID := int64(middleware.GetUserID(c))
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}

	payments, err := h.paymentService.GetGroupPayments(c.Request.Context(), groupID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}
