package user

import (
	"github.com/Longseeyou/DuraFlow/internal/shared/model"
)

type Role int

const (
	ADMIN Role = iota
	MEMBER
)

type User struct {
	model.BaseModel
	Email        string `gorm:"unique"`
	Name         string
	PasswordHash string
	Role         Role
	IsActive     bool
}
