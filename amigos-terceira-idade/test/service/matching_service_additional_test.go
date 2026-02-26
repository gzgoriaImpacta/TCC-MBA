// Package service_test contém os testes adicionais de matching.
package service_test

import (
	"errors"
	"testing"

	"amigos-terceira-idade/internal/domain"
	"amigos-terceira-idade/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestMatchingService_Connect_NotVolunteer testa erro quando não é voluntário.
func TestMatchingService_Connect_NotVolunteer(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	connectionRepo := new(MockConnectionRepository)
	matchingService := service.NewMatchingService(userRepo, connectionRepo)

	elderlyID := uuid.New()
	targetID := uuid.New()

	elderly := &domain.User{ID: elderlyID, UserType: domain.UserTypeElderly}

	connectionRepo.On("Exists", elderlyID, targetID).Return(false, nil)
	userRepo.On("FindByID", elderlyID).Return(elderly, nil)

	// Act
	result, err := matchingService.Connect(elderlyID, targetID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "apenas voluntários podem iniciar conexões", err.Error())
}

// TestMatchingService_GetSuggestions_UserNotFound testa erro quando usuário não existe.
func TestMatchingService_GetSuggestions_UserNotFound(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	connectionRepo := new(MockConnectionRepository)
	matchingService := service.NewMatchingService(userRepo, connectionRepo)

	userID := uuid.New()
	userRepo.On("FindByID", userID).Return(nil, errors.New("usuário não encontrado"))

	// Act
	result, err := matchingService.GetSuggestions(userID, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestMatchingService_GetSuggestions_FilterElderly testa filtro apenas por idosos.
func TestMatchingService_GetSuggestions_FilterElderly(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	connectionRepo := new(MockConnectionRepository)
	matchingService := service.NewMatchingService(userRepo, connectionRepo)

	volunteerID := uuid.New()
	volunteer := &domain.User{
		ID:        volunteerID,
		UserType:  domain.UserTypeVolunteer,
		Interests: []domain.Interest{},
	}

	elderlyID := uuid.New()
	elderly := []domain.User{
		{ID: elderlyID, Name: "Gerson", UserType: domain.UserTypeElderly},
	}

	userRepo.On("FindByID", volunteerID).Return(volunteer, nil)
	userRepo.On("FindByType", domain.UserTypeElderly).Return(elderly, nil)
	connectionRepo.On("Exists", volunteerID, elderlyID).Return(false, nil)

	// Act - filtra apenas idosos
	suggestions, err := matchingService.GetSuggestions(volunteerID, "elderly")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, suggestions, 1)
	assert.Equal(t, "Gerson", suggestions[0].User.Name)
}

// TestMatchingService_GetSuggestions_SkipsConnected testa que pula usuários já conectados.
func TestMatchingService_GetSuggestions_SkipsConnected(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	connectionRepo := new(MockConnectionRepository)
	matchingService := service.NewMatchingService(userRepo, connectionRepo)

	volunteerID := uuid.New()
	volunteer := &domain.User{
		ID:        volunteerID,
		UserType:  domain.UserTypeVolunteer,
		Interests: []domain.Interest{},
	}

	elderlyID1 := uuid.New()
	elderlyID2 := uuid.New()
	elderly := []domain.User{
		{ID: elderlyID1, Name: "Conectado", UserType: domain.UserTypeElderly},
		{ID: elderlyID2, Name: "Não Conectado", UserType: domain.UserTypeElderly},
	}

	userRepo.On("FindByID", volunteerID).Return(volunteer, nil)
	userRepo.On("FindByType", domain.UserTypeElderly).Return(elderly, nil)
	userRepo.On("FindByType", domain.UserTypeInstitution).Return([]domain.User{}, nil)
	connectionRepo.On("Exists", volunteerID, elderlyID1).Return(true, nil)  // Já conectado
	connectionRepo.On("Exists", volunteerID, elderlyID2).Return(false, nil) // Não conectado

	// Act
	suggestions, err := matchingService.GetSuggestions(volunteerID, "")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, suggestions, 1)
	assert.Equal(t, "Não Conectado", suggestions[0].User.Name)
}

// TestMatchingService_GetSuggestions_MultipleMatches testa ordenação por match score.
func TestMatchingService_GetSuggestions_MultipleMatches(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	connectionRepo := new(MockConnectionRepository)
	matchingService := service.NewMatchingService(userRepo, connectionRepo)

	interestID1 := uuid.New()
	interestID2 := uuid.New()

	volunteerID := uuid.New()
	volunteer := &domain.User{
		ID:       volunteerID,
		UserType: domain.UserTypeVolunteer,
		Interests: []domain.Interest{
			{ID: interestID1, Name: "Música"},
			{ID: interestID2, Name: "Xadrez"},
		},
	}

	elderlyID1 := uuid.New()
	elderlyID2 := uuid.New()
	elderly := []domain.User{
		{
			ID:       elderlyID1,
			Name:     "Um Match",
			UserType: domain.UserTypeElderly,
			Interests: []domain.Interest{
				{ID: interestID1, Name: "Música"}, // 1 match
			},
		},
		{
			ID:       elderlyID2,
			Name:     "Dois Matches",
			UserType: domain.UserTypeElderly,
			Interests: []domain.Interest{
				{ID: interestID1, Name: "Música"}, // 2 matches
				{ID: interestID2, Name: "Xadrez"},
			},
		},
	}

	userRepo.On("FindByID", volunteerID).Return(volunteer, nil)
	userRepo.On("FindByType", domain.UserTypeElderly).Return(elderly, nil)
	userRepo.On("FindByType", domain.UserTypeInstitution).Return([]domain.User{}, nil)
	connectionRepo.On("Exists", volunteerID, elderlyID1).Return(false, nil)
	connectionRepo.On("Exists", volunteerID, elderlyID2).Return(false, nil)

	// Act
	suggestions, err := matchingService.GetSuggestions(volunteerID, "")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, suggestions, 2)
	// Deve estar ordenado por score (maior primeiro)
	assert.Equal(t, "Dois Matches", suggestions[0].User.Name)
	assert.Equal(t, 2, suggestions[0].MatchedInterests)
	assert.Equal(t, "Um Match", suggestions[1].User.Name)
	assert.Equal(t, 1, suggestions[1].MatchedInterests)
}

// TestMatchingService_Connect_UserNotFound testa erro quando voluntário não existe.
func TestMatchingService_Connect_UserNotFound(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	connectionRepo := new(MockConnectionRepository)
	matchingService := service.NewMatchingService(userRepo, connectionRepo)

	volunteerID := uuid.New()
	targetID := uuid.New()

	connectionRepo.On("Exists", volunteerID, targetID).Return(false, nil)
	userRepo.On("FindByID", volunteerID).Return(nil, errors.New("usuário não encontrado"))

	// Act
	result, err := matchingService.Connect(volunteerID, targetID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestMatchingService_Connect_TargetNotFound testa erro quando alvo não existe.
func TestMatchingService_Connect_TargetNotFound(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	connectionRepo := new(MockConnectionRepository)
	matchingService := service.NewMatchingService(userRepo, connectionRepo)

	volunteerID := uuid.New()
	targetID := uuid.New()
	volunteer := &domain.User{ID: volunteerID, UserType: domain.UserTypeVolunteer}

	connectionRepo.On("Exists", volunteerID, targetID).Return(false, nil)
	userRepo.On("FindByID", volunteerID).Return(volunteer, nil)
	userRepo.On("FindByID", targetID).Return(nil, errors.New("usuário não encontrado"))

	// Act
	result, err := matchingService.Connect(volunteerID, targetID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestMatchingService_GetConnections_Error testa erro ao buscar conexões.
func TestMatchingService_GetConnections_Error(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	connectionRepo := new(MockConnectionRepository)
	matchingService := service.NewMatchingService(userRepo, connectionRepo)

	userID := uuid.New()
	userRepo.On("FindByID", userID).Return(nil, errors.New("erro"))

	// Act
	result, err := matchingService.GetConnections(userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}
