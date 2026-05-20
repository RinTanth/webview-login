package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"webview-login/backend/auth/model"
	"webview-login/backend/auth/repository"
)

type RegisterInput struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
}

type RegisterService struct {
	repo repository.UserRepo
}

func NewRegisterService(repo repository.UserRepo) *RegisterService {
	return &RegisterService{repo: repo}
}

func (s *RegisterService) Register(in RegisterInput) error {
	if in.Password != in.ConfirmPassword {
		return ErrPasswordMismatch
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.repo.Create(&model.User{
		Username: in.Username,
		Email:    in.Email,
		Password: string(hash),
	})
	if errors.Is(err, repository.ErrDuplicate) {
		return ErrConflict
	}
	return err
}
