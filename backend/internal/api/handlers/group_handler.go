package handlers

import (
	"net/http"
	"strconv"

	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/api/middleware"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/model"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	groupService service.GroupService
}

func NewGroupHandler(gs service.GroupService) *GroupHandler {
	return &GroupHandler{groupService: gs}
}

func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var req model.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// FIX: Use the middleware helper instead of MustGet
	userID := int64(middleware.GetUserID(c))
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user payload"})
		return
	}

	group, err := h.groupService.CreateGroup(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	c.JSON(http.StatusCreated, group)
}

func (h *GroupHandler) AddMember(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var req model.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// FIX: Use the middleware helper instead of MustGet
	requesterID := int64(middleware.GetUserID(c))
	if requesterID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user payload"})
		return
	}

	err = h.groupService.AddMember(c.Request.Context(), groupID, req.UserID, requesterID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added successfully"})
}

func (h *GroupHandler) ListGroups(c *gin.Context) {
	// FIX: Use the middleware helper instead of MustGet
	userID := int64(middleware.GetUserID(c))
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user payload"})
		return
	}

	groups, err := h.groupService.ListGroups(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch groups"})
		return
	}

	c.JSON(http.StatusOK, groups)
}
