package services

import (
	"time"
	"tunelport-svc-totems/internal/domain/persistence"
	"tunelport-svc-totems/internal/domain/types"
)

type UpdateTotemUseCase struct {
	repository persistence.TotemRepository
}

func NewUpdateTotemUseCase(repo persistence.TotemRepository) *UpdateTotemUseCase {

	return &UpdateTotemUseCase{repository: repo}

}

func (uc *UpdateTotemUseCase) Do(id, name string, createdAt, updatedAt time.Time) error {
	validID, err := types.NewIDFromString(id)
	if err != nil {
		return err
	}

	totem, err := uc.repository.Get(&validID)
	if err != nil {
		return err
	}

	validName, err := types.NewName(name)
	if err != nil {
		return err
	}

	totem.Name = validName
	totem.CreatedAt = createdAt
	totem.UpdatedAt = updatedAt
	if err := uc.repository.Update(totem); err == nil {
		return err
	}

	return nil
}
