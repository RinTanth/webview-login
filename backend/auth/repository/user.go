package repository

import (
	"errors"
	"strings"

	"gorm.io/gorm"
	"webview-login/backend/auth/model"
)

var ErrDuplicate = errors.New("duplicate entry")

type UserRepo interface {
	Create(user *model.User) error
	FindByIdentifier(identifier string) (*model.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		// sqlite unique constraint message contains "UNIQUE"
		if isUniqueErr(err) {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

func (r *userRepo) FindByIdentifier(identifier string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ? OR email = ?", identifier, identifier).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func isUniqueErr(err error) bool {
	// SQLite surfaces "UNIQUE constraint failed"; swap for postgres error code "23505"
	return err != nil && strings.Contains(err.Error(), "UNIQUE")
}
