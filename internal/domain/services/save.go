package services

import (
	"time"
	"tunelport-svc-totems/internal/domain/models"
	"tunelport-svc-totems/internal/domain/persistence"
)

type SaveTotemUseCase struct {
	repository persistence.TotemRepository
}

func NewSaveTotemUseCase(repo persistence.TotemRepository) *SaveTotemUseCase {

	return &SaveTotemUseCase{repository: repo}

}

func (uc *SaveTotemUseCase) Do(name string, created_at, updated_at time.Time) (string, error) {
	totem, err := models.NewTotem(name, created_at, updated_at)
	if err != nil {
		return "", err
	}

	err = uc.repository.Save(totem)
	if err != nil {
		return "", err
	}

	return totem.ID.Value(), nil
}
