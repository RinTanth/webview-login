package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"webview-login/backend/auth/service"
)

type LoginHandler struct {
	svc *service.LoginService
}

func NewLoginHandler(svc *service.LoginService) *LoginHandler {
	return &LoginHandler{svc: svc}
}

func (h *LoginHandler) Handle(c *gin.Context) {
	var body struct {
		Identifier string `json:"username" binding:"required"` // accepts username or email
		Password   string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.Login(service.LoginInput{
		Identifier: body.Identifier,
		Password:   body.Password,
	})
	switch {
	case err == nil:
		c.JSON(http.StatusOK, gin.H{
			"token": result.Token,
			"user": gin.H{
				"user_id":  result.UserID,
				"username": result.Username,
			},
		})
	case errors.Is(err, service.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}
