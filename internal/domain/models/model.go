package models

import (
	"time"
	"tunelport-svc-totems/internal/domain/types"
)

type Totem struct {
	ID        types.ID   `json:"id"`
	Name      types.Name `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewTotem(
	name string,
	createdAt time.Time,
	updatedAt time.Time) (*Totem, error) {
	totemName, err := types.NewName(name)
	if err != nil {
		return nil, err
	}

	totemID, err := types.NewID()
	if err != nil {
		return nil, err
	}

	return &Totem{
		ID:        totemID,
		Name:      totemName,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
