package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"webview-login/backend/auth/repository"
	"webview-login/backend/cipher"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type LoginInput struct {
	Identifier string // username or email
	Password   string
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	UserID       uuid.UUID
	Username     string
}

type LoginService struct {
	userRepo    repository.DatabaseRepo
	tokenRepo   repository.RedisRepo
	jwtSecret   string
	emailPepper string
}

func NewLoginService(userRepo repository.DatabaseRepo, tokenRepo repository.RedisRepo, jwtSecret, emailPepper string) *LoginService {
	return &LoginService{
		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		jwtSecret:   jwtSecret,
		emailPepper: emailPepper,
	}
}

func (s *LoginService) Login(ctx context.Context, in LoginInput) (LoginResult, error) {
	hashEmail := cipher.HashEmail(in.Identifier, s.emailPepper)

	user, err := s.userRepo.FindByUsernameOrHashEmail(in.Identifier, hashEmail)
	if err != nil {
		return LoginResult{}, err
	}
	if user == nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	match, err := cipher.VerifyPassword(in.Password, user.HashPassword)
	if err != nil || !match {
		return LoginResult{}, ErrInvalidCredentials
	}

	accessToken, err := s.signJWT(user.UserID, user.Username)
	if err != nil {
		return LoginResult{}, err
	}

	refreshToken, err := newRefreshToken(user.UserID)
	if err != nil {
		return LoginResult{}, err
	}

	if err := s.tokenRepo.SetRefreshToken(ctx, user.UserID.String(), refreshToken); err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.UserID,
		Username:     user.Username,
	}, nil
}

func (s *LoginService) signJWT(userID uuid.UUID, username string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID.String(),
		"username": username,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// newRefreshToken generates a token in the format <userID>.<random_bytes>
// so the server can extract userID from the token without a DB/Redis lookup.
func newRefreshToken(userID uuid.UUID) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s", userID.String(), base64.RawURLEncoding.EncodeToString(b)), nil
}
