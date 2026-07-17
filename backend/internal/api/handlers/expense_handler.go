package handlers

import (
	"net/http"
	"strconv"

	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/api/middleware"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/model"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ExpenseHandler struct {
	expenseService service.ExpenseService
}

func NewExpenseHandler(es service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{expenseService: es}
}

func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var req model.CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := int64(middleware.GetUserID(c))
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user payload"})
		return
	}

	expense, err := h.expenseService.CreateExpense(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, expense)
}

func (h *ExpenseHandler) GetGroupExpenses(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	userID := int64(middleware.GetUserID(c))
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user payload"})
		return
	}

	expenses, err := h.expenseService.GetGroupExpenses(c.Request.Context(), groupID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, expenses)
}

func (h *ExpenseHandler) GetGroupBalances(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	userID := int64(middleware.GetUserID(c))
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user payload"})
		return
	}

	balances, err := h.expenseService.GetGroupBalances(c.Request.Context(), groupID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, balances)
}
