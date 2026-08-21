package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

func (b BaseModel) GetID() uuid.UUID {
	return b.ID
}

func (b BaseModel) GetCreatedAt() *time.Time {
	return &b.CreatedAt
}

func (b BaseModel) GetUpdatedAt() *time.Time {
	return &b.UpdatedAt
}
