package service

import (
	"context"

	"webview-login/backend/auth/repository"
)

type LogoutService struct {
	tokenRepo repository.RedisRepo
}

func NewLogoutService(tokenRepo repository.RedisRepo) *LogoutService {
	return &LogoutService{tokenRepo: tokenRepo}
}

func (s *LogoutService) Logout(ctx context.Context, refreshToken string) error {
	userID, err := parseUserIDFromToken(refreshToken)
	if err != nil {
		return ErrInvalidToken
	}
	return s.tokenRepo.Delete(ctx, userID.String())
}
