package model

import (
	"time"
)

type Credential struct {
	ID             string    `gorm:"column:id;primaryKey;type:char(36)"`
	Username       string    `gorm:"column:username;type:varchar(255);uniqueIndex;not null"`
	HashedPassword string    `gorm:"column:hashed_password;type:varchar(255);not null"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
