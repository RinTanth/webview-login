package service

import (
	"context"
	"log/slog"

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
		slog.Warn("logout failed: invalid token format")
		return ErrInvalidToken
	}

	if err := s.tokenRepo.Delete(ctx, userID.String()); err != nil {
		slog.Error("logout failed: redis error", "user_id", userID, "error", err)
		return err
	}

	slog.Info("logout success", "user_id", userID)
	return nil
}
