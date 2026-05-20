package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"webview-login/backend/auth/service"
)

type LogoutHandler struct {
	svc *service.LogoutService
}

func NewLogoutHandler(svc *service.LogoutService) *LogoutHandler {
	return &LogoutHandler{svc: svc}
}

func (h *LogoutHandler) Handle(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.Logout(c.Request.Context(), body.RefreshToken)
	switch {
	case err == nil:
		c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
	case errors.Is(err, service.ErrInvalidToken):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}
