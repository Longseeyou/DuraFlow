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
	email 			string	`gorm:"unique"`
	name 			string 
	password_hash 	string
	role 			role
	isActive 		bool

}