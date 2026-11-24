package rest

import (
	"encoding/json"
	"net/http"
	"tunelport-svc-totems/internal/domain/persistence"
	"tunelport-svc-totems/internal/domain/services"
)

type RESTSaveUseCase struct {
	repository persistence.TotemRepository
}

func NewRESTSaveUseCase(repository persistence.TotemRepository) *RESTSaveUseCase {
	return &RESTSaveUseCase{
		repository: repository,
	}
}

func (uc *RESTSaveUseCase) Save(w http.ResponseWriter, r *http.Request) {
	var dto totemDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	saveUC := services.NewSaveTotemUseCase(uc.repository)
	id, err := saveUC.Do(
		dto.Name,
		dto.CreatedAt,
		dto.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto)
}
