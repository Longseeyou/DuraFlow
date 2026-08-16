package user

import (
	"github.com/Longseeyou/DuraFlow/internal/shared/model"
)

type Role string

const (
	ADMIN Role = "ADMIN"
	USER  Role = "USER"
)

type User struct {
	model.BaseModel
	Email        string `gorm:"uniqueIndex; not null"`
	Name         string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`
	Role         Role
	IsActive     bool
}
