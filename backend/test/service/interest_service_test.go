// Package service_test contém os testes dos serviços da aplicação.
package service_test

import (
	"errors"
	"testing"

	"amigos-terceira-idade/internal/domain"
	"amigos-terceira-idade/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestInterestService_GetAll_Success testa listagem de todos os interesses.
func TestInterestService_GetAll_Success(t *testing.T) {
	// Arrange
	interestRepo := new(MockInterestRepository)
	interestService := service.NewInterestService(interestRepo)

	interests := []domain.Interest{
		{ID: uuid.New(), Name: "Música", Icon: "🎵"},
		{ID: uuid.New(), Name: "Xadrez", Icon: "♟️"},
		{ID: uuid.New(), Name: "Leitura", Icon: "📚"},
	}

	interestRepo.On("FindAll").Return(interests, nil)

	// Act
	result, err := interestService.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "Música", result[0].Name)
	interestRepo.AssertExpectations(t)
}

// TestInterestService_GetAll_Empty testa listagem quando não há interesses.
func TestInterestService_GetAll_Empty(t *testing.T) {
	// Arrange
	interestRepo := new(MockInterestRepository)
	interestService := service.NewInterestService(interestRepo)

	interestRepo.On("FindAll").Return([]domain.Interest{}, nil)

	// Act
	result, err := interestService.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

// TestInterestService_GetByID_Success testa busca de interesse por ID.
func TestInterestService_GetByID_Success(t *testing.T) {
	// Arrange
	interestRepo := new(MockInterestRepository)
	interestService := service.NewInterestService(interestRepo)

	interestID := uuid.New()
	interest := &domain.Interest{
		ID:   interestID,
		Name: "Caminhadas",
		Icon: "🚶",
	}

	interestRepo.On("FindByID", interestID).Return(interest, nil)

	// Act
	result, err := interestService.GetByID(interestID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Caminhadas", result.Name)
}

// TestInterestService_GetByID_NotFound testa erro quando interesse não existe.
func TestInterestService_GetByID_NotFound(t *testing.T) {
	// Arrange
	interestRepo := new(MockInterestRepository)
	interestService := service.NewInterestService(interestRepo)

	interestID := uuid.New()
	interestRepo.On("FindByID", interestID).Return(nil, errors.New("interesse não encontrado"))

	// Act
	result, err := interestService.GetByID(interestID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestInterestService_SeedDefaults_Success testa inserção de interesses padrão.
func TestInterestService_SeedDefaults_Success(t *testing.T) {
	// Arrange
	interestRepo := new(MockInterestRepository)
	interestService := service.NewInterestService(interestRepo)

	interestRepo.On("SeedDefaults").Return(nil)

	// Act
	err := interestService.SeedDefaults()

	// Assert
	assert.NoError(t, err)
	interestRepo.AssertExpectations(t)
}

// TestInterestService_SeedDefaults_Error testa erro ao inserir interesses padrão.
func TestInterestService_SeedDefaults_Error(t *testing.T) {
	// Arrange
	interestRepo := new(MockInterestRepository)
	interestService := service.NewInterestService(interestRepo)

	interestRepo.On("SeedDefaults").Return(errors.New("erro no banco"))

	// Act
	err := interestService.SeedDefaults()

	// Assert
	assert.Error(t, err)
}
