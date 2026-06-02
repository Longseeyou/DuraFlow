package dto

import (
	"time"

	"github.com/google/uuid"
)

type ResponseDto struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BaseModel interface {
	GetID() uuid.UUID
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

func (rD *ResponseDto) BaseModelToResponseDto(bM BaseModel) {
	rD.ID = bM.GetID()
	rD.CreatedAt = bM.GetCreatedAt()
	rD.UpdatedAt = bM.GetUpdatedAt()
}
