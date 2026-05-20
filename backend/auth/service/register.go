package service

import (
	"errors"
	"log/slog"

	"webview-login/backend/auth/model"
	"webview-login/backend/auth/repository"
	"webview-login/backend/cipher"
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
		slog.Warn("register failed: password mismatch", "username", in.Username)
		return ErrPasswordMismatch
	}

	hashEmail := cipher.HashEmail(in.Email, s.emailPepper)

	encEmail, err := cipher.EncryptEmail(in.Email, s.aesKey)
	if err != nil {
		slog.Error("register failed: email encryption error", "username", in.Username, "error", err)
		return err
	}

	hashPassword, err := cipher.HashPassword(in.Password)
	if err != nil {
		slog.Error("register failed: password hashing error", "username", in.Username, "error", err)
		return err
	}

	err = s.repo.Create(&model.UserInfo{
		Username:     in.Username,
		HashEmail:    hashEmail,
		EncEmail:     encEmail,
		HashPassword: hashPassword,
	})
	if errors.Is(err, repository.ErrDuplicate) {
		slog.Warn("register failed: duplicate username or email", "username", in.Username)
		return ErrConflict
	}
	if err != nil {
		slog.Error("register failed: database error", "username", in.Username, "error", err)
		return err
	}

	slog.Info("register success", "username", in.Username)
	return nil
}
