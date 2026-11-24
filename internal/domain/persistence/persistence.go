package persistence

import (
	"tunelport-svc-totems/internal/domain/models"
	"tunelport-svc-totems/internal/domain/types"
)

type TotemRepository interface {
	Save(totem *models.Totem) error
	Get(id *types.ID) (*models.Totem, error)
	Delete(id *types.ID) error
	Update(totem *models.Totem) error
	List() ([]*models.Totem, error)
}
