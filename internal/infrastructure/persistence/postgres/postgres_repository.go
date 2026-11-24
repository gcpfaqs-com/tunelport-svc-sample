package postgres

import (
	"errors"
	"time"
	"tunelport-svc-totems/internal/domain/models"
	"tunelport-svc-totems/internal/domain/persistence"
	"tunelport-svc-totems/internal/domain/types"

	"gorm.io/gorm"
)

type PostgresTotemRepository struct {
	db *gorm.DB
}

type totemRecord struct {
	ID        string    `gorm:"primaryKey;type:uuid;not null"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (totemRecord) TableName() string {
	return "totems"
}

func NewPostgresTotemRepository(db *gorm.DB) persistence.TotemRepository {
	if db == nil {
		panic("postgres repository requires a valid gorm db instance")
	}

	if err := db.AutoMigrate(&totemRecord{}); err != nil {
		panic(err)
	}

	return &PostgresTotemRepository{db: db}
}

func (r *PostgresTotemRepository) Save(totem *models.Totem) error {
	if totem == nil {
		return errors.New("totem is nil")
	}
	if r.db == nil {
		return errors.New("database connection is nil")
	}
	if totem.ID.Value() == "" ||
		totem.Name.Value() == "" ||
		totem.CreatedAt.IsZero() ||
		totem.UpdatedAt.IsZero() {
		return errors.New("totem is incomplete")
	}

	payload := totemRecord{
		ID:        totem.ID.Value(),
		Name:      totem.Name.Value(),
		CreatedAt: totem.CreatedAt,
		UpdatedAt: totem.UpdatedAt,
	}

	return r.db.Save(&payload).Error
}

func (r *PostgresTotemRepository) Get(id *types.ID) (*models.Totem, error) {
	if id == nil {
		return nil, errors.New("totem id is nil")
	}
	if r.db == nil {
		return nil, errors.New("database connection is nil")
	}

	var payload totemRecord
	result := r.db.Where("id = ?", id.Value()).First(&payload)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}

	totemID, err := types.NewIDFromString(payload.ID)
	if err != nil {
		return nil, err
	}

	totemName, err := types.NewName(payload.Name)
	if err != nil {
		return nil, err
	}

	return &models.Totem{
		ID:        totemID,
		Name:      totemName,
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.UpdatedAt,
	}, nil
}

func (r *PostgresTotemRepository) Delete(id *types.ID) error {
	if id == nil {
		return errors.New("totem id is nil")
	}
	if r.db == nil {
		return errors.New("database connection is nil")
	}

	result := r.db.Delete(&totemRecord{}, "id = ?", id.Value())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *PostgresTotemRepository) Update(totem *models.Totem) error {
	if totem == nil {
		return errors.New("totem is nil")
	}
	if r.db == nil {
		return errors.New("database connection is nil")
	}
	if totem.ID.Value() == "" ||
		totem.Name.Value() == "" ||
		totem.CreatedAt.IsZero() ||
		totem.UpdatedAt.IsZero() {
		return errors.New("totem is incomplete")
	}

	payload := totemRecord{
		ID:        totem.ID.Value(),
		Name:      totem.Name.Value(),
		CreatedAt: totem.CreatedAt,
		UpdatedAt: totem.UpdatedAt,
	}

	result := r.db.Model(&totemRecord{}).Where("id = ?", payload.ID).Updates(payload)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *PostgresTotemRepository) List() ([]*models.Totem, error) {
	if r.db == nil {
		return nil, errors.New("database connection is nil")
	}

	var payloads []totemRecord
	if err := r.db.Find(&payloads).Error; err != nil {
		return nil, err
	}

	totems := make([]*models.Totem, 0, len(payloads))
	for _, payload := range payloads {
		totemID, err := types.NewIDFromString(payload.ID)
		if err != nil {
			return nil, err
		}

		totemName, err := types.NewName(payload.Name)
		if err != nil {
			return nil, err
		}

		totems = append(totems, &models.Totem{
			ID:        totemID,
			Name:      totemName,
			CreatedAt: payload.CreatedAt,
			UpdatedAt: payload.UpdatedAt,
		})
	}

	return totems, nil
}
