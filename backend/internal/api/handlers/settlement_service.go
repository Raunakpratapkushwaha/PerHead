package handlers

import (
	"net/http"
	"strconv"

	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/api/middleware"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type SettlementHandler struct {
	settlementService service.SettlementService
}

func NewSettlementHandler(ss service.SettlementService) *SettlementHandler {
	return &SettlementHandler{settlementService: ss}
}

func (h *SettlementHandler) GetGroupSettlements(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// FIX: Safely extract the user ID using the middleware helper
	userID := int64(middleware.GetUserID(c))
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user payload"})
		return
	}

	settlements, err := h.settlementService.GetSimplifiedDebts(c.Request.Context(), groupID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"group_id":    groupID,
		"settlements": settlements,
	})
}