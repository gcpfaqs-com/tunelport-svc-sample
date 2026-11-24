package rest

import (
	"encoding/json"
	"net/http"
	"tunelport-svc-totems/internal/domain/persistence"
	"tunelport-svc-totems/internal/domain/services"
)

type RESTUpdateUseCase struct {
	repository persistence.TotemRepository
}

func NewRESTUpdateUseCase(repository persistence.TotemRepository) *RESTUpdateUseCase {
	return &RESTUpdateUseCase{
		repository: repository,
	}
}

func (uc *RESTUpdateUseCase) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	var dto totemDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updateUC := services.NewUpdateTotemUseCase(uc.repository)
	err := updateUC.Do(
		idStr,
		dto.Name,
		dto.CreatedAt,
		dto.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto.ID = idStr
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto)
}
