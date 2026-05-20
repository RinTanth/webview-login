package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"webview-login/backend/cipher"
	"webview-login/backend/auth/repository"
)

type LoginInput struct {
	Identifier string // username or email
	Password   string
}

type LoginResult struct {
	Token    string
	UserID   uuid.UUID
	Username string
}

type LoginService struct {
	repo        repository.UserRepo
	jwtSecret   string
	emailPepper string
}

func NewLoginService(repo repository.UserRepo, jwtSecret, emailPepper string) *LoginService {
	return &LoginService{repo: repo, jwtSecret: jwtSecret, emailPepper: emailPepper}
}

func (s *LoginService) Login(in LoginInput) (LoginResult, error) {
	// hash the identifier so we can match against hash_email if it's an email
	hashEmail := cipher.HashEmail(in.Identifier, s.emailPepper)

	user, err := s.repo.FindByUsernameOrHashEmail(in.Identifier, hashEmail)
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

	token, err := s.signJWT(user.UserID, user.Username)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		Token:    token,
		UserID:   user.UserID,
		Username: user.Username,
	}, nil
}

func (s *LoginService) signJWT(userID uuid.UUID, username string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID.String(),
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

