package rest

import (
	"encoding/json"
	"net/http"
	"tunelport-svc-totems/internal/domain/persistence"
	"tunelport-svc-totems/internal/domain/services"
)

type RESTListUseCase struct {
	repository persistence.TotemRepository
}

func NewRESTListUseCase(repository persistence.TotemRepository) *RESTListUseCase {
	return &RESTListUseCase{
		repository: repository,
	}
}

func (uc *RESTListUseCase) List(w http.ResponseWriter, r *http.Request) {
	listUC := services.NewListTotemUseCase(uc.repository)
	totems, err := listUC.Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := make([]totemDTO, 0, len(totems))
	for _, a := range totems {
		response = append(response, totemDTO{
			ID:        a.ID.Value(),
			Name:      a.Name.Value(),
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
