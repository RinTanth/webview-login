package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"webview-login/backend/auth/service"
)

type RegisterHandler struct {
	svc *service.RegisterService
}

func NewRegisterHandler(svc *service.RegisterService) *RegisterHandler {
	return &RegisterHandler{svc: svc}
}

func (h *RegisterHandler) Handle(c *gin.Context) {
	var body struct {
		Username        string `json:"username"        binding:"required"`
		Email           string `json:"email"           binding:"required,email"`
		Password        string `json:"password"        binding:"required,min=6"`
		ConfirmPassword string `json:"confirmPassword" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.Register(service.RegisterInput{
		Username:        body.Username,
		Email:           body.Email,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
	})
	switch {
	case err == nil:
		c.JSON(http.StatusCreated, gin.H{"message": "registered successfully"})
	case errors.Is(err, service.ErrPasswordMismatch):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}
