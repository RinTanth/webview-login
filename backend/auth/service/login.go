package service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"webview-login/backend/auth/repository"
)

type LoginInput struct {
	Identifier string // username or email
	Password   string
}

type LoginResult struct {
	Token    string
	ID       uint
	Username string
	Email    string
}

type LoginService struct {
	repo repository.UserRepo
}

func NewLoginService(repo repository.UserRepo) *LoginService {
	return &LoginService{repo: repo}
}

func (s *LoginService) Login(in LoginInput) (LoginResult, error) {
	user, err := s.repo.FindByIdentifier(in.Identifier)
	if err != nil {
		return LoginResult{}, err
	}
	if user == nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)); err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	token, err := signJWT(user.ID, user.Username)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		Token:    token,
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func signJWT(id uint, username string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      id,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
