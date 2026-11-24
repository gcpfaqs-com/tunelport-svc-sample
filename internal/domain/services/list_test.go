package services_test

import (
	"testing"
	"tunelport-svc-totems/internal/domain/models"
	"tunelport-svc-totems/internal/domain/services"
	"tunelport-svc-totems/internal/domain/test"

	"github.com/stretchr/testify/assert"
)

func TestItShouldResultAList(t *testing.T) {
	repository := test.NewStubRepository()

	sut := services.NewListTotemUseCase(repository)
	list, err := sut.Do()

	assert.Nil(t, err)
	assert.IsType(t, []*models.Totem{}, list)
	assert.True(t, repository.ListHaveBeenCalled(), sut)
}
