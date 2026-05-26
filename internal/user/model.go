package user 

import (
	"time"
	"github.com/Longseeyou/DuraFlow/internal/shared/model"
)

type role int

const (
	admin role = iota
	member
)

type User struct {
	model.BaseModel
	Email 			string	`gorm:"unique"`
	Name 			string 
	PasswordHash 	string
	Role 			role
	IsActive 		bool

}