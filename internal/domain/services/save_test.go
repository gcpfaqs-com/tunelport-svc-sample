package services_test

import (
	"time"

	"testing"
	"tunelport-svc-totems/internal/domain/services"
	"tunelport-svc-totems/internal/domain/test"
	"tunelport-svc-totems/internal/domain/types"

	"github.com/stretchr/testify/assert"
)

var newName string = "any valid name"

func TestItShouldSaveAResource(t *testing.T) {
	repository := test.NewStubRepository()

	sut := services.NewSaveTotemUseCase(repository)
	_, err := sut.Do(newName, time.Now(), time.Now())

	assert.Nil(t, err)
	assert.True(t, repository.SaveHaveBeenCalled(), sut)
}

func TestItShouldResultAnErrorOnInvalidName(t *testing.T) {
	newName = ""
	repository := test.NewStubRepository()

	sut := services.NewSaveTotemUseCase(repository)
	_, err := sut.Do(newName, time.Now(), time.Now())

	assert.NotNil(t, err)
	assert.IsType(t, types.ErrInvalidName, err)
}
