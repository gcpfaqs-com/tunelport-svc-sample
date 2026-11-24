package services

import (
	"tunelport-svc-totems/internal/domain/models"
	"tunelport-svc-totems/internal/domain/persistence"
)

type ListTotemUseCase struct {
	repository persistence.TotemRepository
}

func NewListTotemUseCase(repo persistence.TotemRepository) *ListTotemUseCase {
	return &ListTotemUseCase{repository: repo}
}

func (uc *ListTotemUseCase) Do() ([]*models.Totem, error) {
	totems, err := uc.repository.List()
	if err != nil {
		return nil, err
	}

	return totems, nil
}
