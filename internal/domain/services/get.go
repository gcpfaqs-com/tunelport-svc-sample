package services

import (
	"tunelport-svc-totems/internal/domain/models"
	"tunelport-svc-totems/internal/domain/persistence"
	"tunelport-svc-totems/internal/domain/types"
)

type GetTotemUseCase struct {
	repository persistence.TotemRepository
}

func NewGetTotemUseCase(repo persistence.TotemRepository) *GetTotemUseCase {
	return &GetTotemUseCase{repository: repo}
}

func (uc *GetTotemUseCase) Do(id string) (*models.Totem, error) {
	validID, err := types.NewIDFromString(id)
	if err != nil {
		return nil, err
	}

	totem, err := uc.repository.Get(&validID)
	if err != nil {
		return nil, err
	}

	return totem, nil
}
