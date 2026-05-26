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
	email 			String	`gorm:"unique"`
	name 			String 
	password_hash 	String
	role 			Role
	isActive 		Bool

}