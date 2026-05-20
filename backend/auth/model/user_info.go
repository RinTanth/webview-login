package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserInfo struct {
	UserID       uuid.UUID `gorm:"type:uuid;primaryKey"                json:"user_id"`
	Username     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	HashEmail    string    `gorm:"type:varchar;uniqueIndex;not null"   json:"-"`
	EncEmail     string    `gorm:"type:varchar;not null"               json:"-"`
	HashPassword string    `gorm:"type:varchar;not null"               json:"-"`
	CreatedAt    time.Time `                                           json:"created_at"`
	UpdatedAt    time.Time `                                           json:"updated_at"`
}

func (UserInfo) TableName() string {
	return "user_info"
}

func (u *UserInfo) BeforeCreate(tx *gorm.DB) error {
	u.UserID = uuid.New()
	return nil
}
