// Package service_test contém os testes adicionais para aumentar a cobertura.
package service_test

import (
	"testing"

	"amigos-terceira-idade/internal/domain"
	"amigos-terceira-idade/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestAuthService_RefreshToken_Success testa renovação de token com sucesso.
func TestAuthService_RefreshToken_Success(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	authService := service.NewAuthService(userRepo, interestRepo, getTestJWTConfig())

	// Primeiro registra um usuário para obter um refresh token válido
	req := service.RegisterRequest{
		Name:     "Refresh User",
		Email:    "refresh@email.com",
		Password: "senha123",
		UserType: "VOLUNTEER",
	}

	userRepo.On("ExistsByEmail", req.Email).Return(false, nil)
	userRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	result, _ := authService.Register(req)

	// Prepara o mock para FindByID
	userRepo.On("FindByID", mock.AnythingOfType("uuid.UUID")).Return(&domain.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Name:     req.Name,
		UserType: domain.UserTypeVolunteer,
	}, nil)

	// Act
	newTokens, err := authService.RefreshToken(result.RefreshToken)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, newTokens)
	assert.NotEmpty(t, newTokens.AccessToken)
	assert.NotEmpty(t, newTokens.RefreshToken)
}

// TestAuthService_RefreshToken_InvalidToken testa renovação com token inválido.
func TestAuthService_RefreshToken_InvalidToken(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	authService := service.NewAuthService(userRepo, interestRepo, getTestJWTConfig())

	// Act
	result, err := authService.RefreshToken("token-invalido")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestAuthService_Register_DBError testa erro ao salvar no banco.
func TestAuthService_Register_DBError(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	authService := service.NewAuthService(userRepo, interestRepo, getTestJWTConfig())

	req := service.RegisterRequest{
		Name:     "Erro User",
		Email:    "erro@email.com",
		Password: "senha123",
		UserType: "VOLUNTEER",
	}

	userRepo.On("ExistsByEmail", req.Email).Return(false, nil)
	userRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(assert.AnError)

	// Act
	result, err := authService.Register(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestAuthService_Register_InterestFetchError testa erro ao buscar interesses.
func TestAuthService_Register_InterestFetchError(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	authService := service.NewAuthService(userRepo, interestRepo, getTestJWTConfig())

	interestID := uuid.New()
	req := service.RegisterRequest{
		Name:        "Interest Error",
		Email:       "interesterror@email.com",
		Password:    "senha123",
		UserType:    "VOLUNTEER",
		InterestIDs: []uuid.UUID{interestID},
	}

	userRepo.On("ExistsByEmail", req.Email).Return(false, nil)
	interestRepo.On("FindByIDs", req.InterestIDs).Return([]domain.Interest{}, assert.AnError)

	// Act
	result, err := authService.Register(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}
