package service

import (
	"errors"

	"webview-login/backend/cipher"
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
	repo        repository.DatabaseRepo
	emailPepper string
	aesKey      string
}

func NewRegisterService(repo repository.DatabaseRepo, emailPepper, aesKey string) *RegisterService {
	return &RegisterService{repo: repo, emailPepper: emailPepper, aesKey: aesKey}
}

func (s *RegisterService) Register(in RegisterInput) error {
	if in.Password != in.ConfirmPassword {
		return ErrPasswordMismatch
	}

	hashEmail := cipher.HashEmail(in.Email, s.emailPepper)

	encEmail, err := cipher.EncryptEmail(in.Email, s.aesKey)
	if err != nil {
		return err
	}

	hashPassword, err := cipher.HashPassword(in.Password)
	if err != nil {
		return err
	}

	err = s.repo.Create(&model.UserInfo{
		Username:     in.Username,
		HashEmail:    hashEmail,
		EncEmail:     encEmail,
		HashPassword: hashPassword,
	})
	if errors.Is(err, repository.ErrDuplicate) {
		return ErrConflict
	}
	return err
}
