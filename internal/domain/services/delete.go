package services

import (
	"tunelport-svc-totems/internal/domain/persistence"
	"tunelport-svc-totems/internal/domain/types"
)

type DeleteTotemUseCase struct {
	repository persistence.TotemRepository
}

func NewDeleteTotemUseCase(repo persistence.TotemRepository) *DeleteTotemUseCase {
	return &DeleteTotemUseCase{repository: repo}
}

func (uc *DeleteTotemUseCase) Do(id string) error {
	validID, err := types.NewIDFromString(id)
	if err != nil {
		return err
	}

	err = uc.repository.Delete(&validID)
	if err != nil {
		return err
	}

	return nil
}
