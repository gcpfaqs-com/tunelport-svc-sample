package rest

import (
	"encoding/json"
	"net/http"
	"tunelport-svc-totems/internal/domain/persistence"
	"tunelport-svc-totems/internal/domain/services"
)

type RESTGetUseCase struct {
	repository persistence.TotemRepository
}

func NewRESTGetUseCase(repository persistence.TotemRepository) *RESTGetUseCase {
	return &RESTGetUseCase{
		repository: repository,
	}
}

func (uc *RESTGetUseCase) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	getUC := services.NewGetTotemUseCase(uc.repository)
	totem, err := getUC.Do(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := totemDTO{
		ID:        totem.ID.Value(),
		Name:      totem.Name.Value(),
		CreatedAt: totem.CreatedAt,
		UpdatedAt: totem.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
