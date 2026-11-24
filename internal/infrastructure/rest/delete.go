package rest

import (
	"net/http"
	"tunelport-svc-totems/internal/domain/persistence"
	"tunelport-svc-totems/internal/domain/services"
)

type RESTDeleteUseCase struct {
	repository persistence.TotemRepository
}

func NewRESTDeleteUseCase(repository persistence.TotemRepository) *RESTDeleteUseCase {
	return &RESTDeleteUseCase{
		repository: repository,
	}
}

func (uc *RESTDeleteUseCase) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	deleteUC := services.NewDeleteTotemUseCase(uc.repository)
	if err := deleteUC.Do(idStr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
