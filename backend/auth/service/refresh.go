package service

import (
	"context"
	"strings"
	"time"

	"webview-login/backend/auth/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RefreshResult struct {
	AccessToken  string
	RefreshToken string
}

type RefreshService struct {
	userRepo  repository.DatabaseRepo
	tokenRepo repository.RedisRepo
	jwtSecret string
}

func NewRefreshService(userRepo repository.DatabaseRepo, tokenRepo repository.RedisRepo, jwtSecret string) *RefreshService {
	return &RefreshService{userRepo: userRepo, tokenRepo: tokenRepo, jwtSecret: jwtSecret}
}

func (s *RefreshService) Refresh(ctx context.Context, refreshToken string) (RefreshResult, error) {
	userID, err := parseUserIDFromToken(refreshToken)
	if err != nil {
		return RefreshResult{}, ErrInvalidToken
	}

	valid, err := s.tokenRepo.GetRefreshToken(ctx, userID.String(), refreshToken)
	if err != nil {
		return RefreshResult{}, err
	}
	if !valid {
		return RefreshResult{}, ErrInvalidToken
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return RefreshResult{}, ErrInvalidToken
	}

	accessToken, err := signAccessJWT(userID, user.Username, s.jwtSecret)
	if err != nil {
		return RefreshResult{}, err
	}

	newRefresh, err := newRefreshToken(userID)
	if err != nil {
		return RefreshResult{}, err
	}

	if err := s.tokenRepo.SetRefreshToken(ctx, userID.String(), newRefresh); err != nil {
		return RefreshResult{}, err
	}

	return RefreshResult{AccessToken: accessToken, RefreshToken: newRefresh}, nil
}

func parseUserIDFromToken(token string) (uuid.UUID, error) {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return uuid.UUID{}, ErrInvalidToken
	}
	return uuid.Parse(parts[0])
}

func signAccessJWT(userID uuid.UUID, username, secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID.String(),
		"username": username,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
