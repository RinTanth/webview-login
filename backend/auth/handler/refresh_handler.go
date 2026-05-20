package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"webview-login/backend/auth/service"
)

type RefreshHandler struct {
	svc *service.RefreshService
}

func NewRefreshHandler(svc *service.RefreshService) *RefreshHandler {
	return &RefreshHandler{svc: svc}
}

func (h *RefreshHandler) Handle(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.Refresh(c.Request.Context(), body.RefreshToken)
	switch {
	case err == nil:
		c.JSON(http.StatusOK, gin.H{
			"access_token":  result.AccessToken,
			"refresh_token": result.RefreshToken,
		})
	case errors.Is(err, service.ErrInvalidToken):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}
