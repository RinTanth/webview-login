package repository

import (
	"errors"
	"strings"

	"gorm.io/gorm"
	"webview-login/backend/auth/model"
)

var ErrDuplicate = errors.New("duplicate entry")

type UserRepo interface {
	Create(user *model.UserInfo) error
	FindByUsernameOrHashEmail(identifier, hashEmail string) (*model.UserInfo, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *model.UserInfo) error {
	if err := r.db.Create(user).Error; err != nil {
		if isUniqueErr(err) {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

// FindByUsernameOrHashEmail looks up by plain username OR pre-hashed email.
// The caller is responsible for hashing the email before passing it in.
func (r *userRepo) FindByUsernameOrHashEmail(identifier, hashEmail string) (*model.UserInfo, error) {
	var user model.UserInfo
	err := r.db.
		Where("username = ? OR hash_email = ?", identifier, hashEmail).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func isUniqueErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "23505")
}
