// Package service_test contém os testes dos serviços da aplicação.
package service_test

import (
	"errors"
	"testing"

	"amigos-terceira-idade/internal/domain"
	"amigos-terceira-idade/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestUserService_GetByID_Success testa busca de usuário pelo ID com sucesso.
func TestUserService_GetByID_Success(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	userService := service.NewUserService(userRepo, interestRepo)

	userID := uuid.New()
	expectedUser := &domain.User{
		ID:    userID,
		Name:  "João Silva",
		Email: "joao@email.com",
	}

	userRepo.On("FindByID", userID).Return(expectedUser, nil)

	// Act
	result, err := userService.GetByID(userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.Name, result.Name)
	userRepo.AssertExpectations(t)
}

// TestUserService_GetByID_NotFound testa erro quando usuário não existe.
func TestUserService_GetByID_NotFound(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	userService := service.NewUserService(userRepo, interestRepo)

	userID := uuid.New()
	userRepo.On("FindByID", userID).Return(nil, errors.New("usuário não encontrado"))

	// Act
	result, err := userService.GetByID(userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestUserService_UpdateProfile_Success testa atualização de perfil com sucesso.
func TestUserService_UpdateProfile_Success(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	userService := service.NewUserService(userRepo, interestRepo)

	userID := uuid.New()
	existingUser := &domain.User{
		ID:    userID,
		Name:  "João Silva",
		Email: "joao@email.com",
		Age:   30,
	}

	updatedUser := &domain.User{
		ID:    userID,
		Name:  "João da Silva",
		Email: "joao@email.com",
		Age:   31,
		Bio:   "Nova bio",
	}

	req := service.UpdateProfileRequest{
		Name: "João da Silva",
		Age:  31,
		Bio:  "Nova bio",
	}

	userRepo.On("FindByID", userID).Return(existingUser, nil).Once()
	userRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)
	userRepo.On("FindByID", userID).Return(updatedUser, nil).Once()

	// Act
	result, err := userService.UpdateProfile(userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "João da Silva", result.Name)
	assert.Equal(t, 31, result.Age)
}

// TestUserService_UpdateProfile_WithInterests testa atualização com interesses.
func TestUserService_UpdateProfile_WithInterests(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	userService := service.NewUserService(userRepo, interestRepo)

	userID := uuid.New()
	interestID := uuid.New()
	existingUser := &domain.User{
		ID:   userID,
		Name: "Maria",
	}

	interests := []domain.Interest{
		{ID: interestID, Name: "Música"},
	}

	req := service.UpdateProfileRequest{
		InterestIDs: []uuid.UUID{interestID},
	}

	userRepo.On("FindByID", userID).Return(existingUser, nil).Once()
	interestRepo.On("FindByIDs", req.InterestIDs).Return(interests, nil)
	userRepo.On("UpdateInterests", userID, interests).Return(nil)
	userRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)
	userRepo.On("FindByID", userID).Return(existingUser, nil).Once()

	// Act
	result, err := userService.UpdateProfile(userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	interestRepo.AssertExpectations(t)
}

// TestUserService_UpdateProfile_UserNotFound testa erro quando usuário não existe.
func TestUserService_UpdateProfile_UserNotFound(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	userService := service.NewUserService(userRepo, interestRepo)

	userID := uuid.New()
	req := service.UpdateProfileRequest{Name: "Novo Nome"}

	userRepo.On("FindByID", userID).Return(nil, errors.New("usuário não encontrado"))

	// Act
	result, err := userService.UpdateProfile(userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestUserService_ListByType_Success testa listagem por tipo com sucesso.
func TestUserService_ListByType_Success(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	userService := service.NewUserService(userRepo, interestRepo)

	volunteers := []domain.User{
		{ID: uuid.New(), Name: "Ricardo", UserType: domain.UserTypeVolunteer},
		{ID: uuid.New(), Name: "Ana", UserType: domain.UserTypeVolunteer},
	}

	userRepo.On("FindByType", domain.UserTypeVolunteer).Return(volunteers, nil)

	// Act
	result, err := userService.ListByType(domain.UserTypeVolunteer)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

// TestUserService_Deactivate_Success testa desativação de usuário com sucesso.
func TestUserService_Deactivate_Success(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	userService := service.NewUserService(userRepo, interestRepo)

	userID := uuid.New()
	user := &domain.User{
		ID:       userID,
		Name:     "João",
		IsActive: true,
	}

	userRepo.On("FindByID", userID).Return(user, nil)
	userRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)

	// Act
	err := userService.Deactivate(userID)

	// Assert
	assert.NoError(t, err)
	assert.False(t, user.IsActive) // Deve estar desativado
	userRepo.AssertExpectations(t)
}

// TestUserService_Deactivate_UserNotFound testa erro quando usuário não existe.
func TestUserService_Deactivate_UserNotFound(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	userService := service.NewUserService(userRepo, interestRepo)

	userID := uuid.New()
	userRepo.On("FindByID", userID).Return(nil, errors.New("usuário não encontrado"))

	// Act
	err := userService.Deactivate(userID)

	// Assert
	assert.Error(t, err)
}
