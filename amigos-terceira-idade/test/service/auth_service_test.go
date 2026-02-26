// Package service_test contém os testes dos serviços da aplicação.
package service_test

import (
	"errors"
	"testing"

	"amigos-terceira-idade/internal/config"
	"amigos-terceira-idade/internal/domain"
	"amigos-terceira-idade/internal/repository"
	"amigos-terceira-idade/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository implementa repository.UserRepositoryInterface para testes.
type MockUserRepository struct {
	mock.Mock
}

// Garante que implementa a interface
var _ repository.UserRepositoryInterface = (*MockUserRepository)(nil)

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) FindByType(userType domain.UserType) ([]domain.User, error) {
	args := m.Called(userType)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) AddInterests(userID uuid.UUID, interests []domain.Interest) error {
	args := m.Called(userID, interests)
	return args.Error(0)
}

func (m *MockUserRepository) RemoveInterest(userID uuid.UUID, interestID uuid.UUID) error {
	args := m.Called(userID, interestID)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateInterests(userID uuid.UUID, interests []domain.Interest) error {
	args := m.Called(userID, interests)
	return args.Error(0)
}

// MockInterestRepository implementa repository.InterestRepositoryInterface para testes.
type MockInterestRepository struct {
	mock.Mock
}

// Garante que implementa a interface
var _ repository.InterestRepositoryInterface = (*MockInterestRepository)(nil)

func (m *MockInterestRepository) Create(interest *domain.Interest) error {
	args := m.Called(interest)
	return args.Error(0)
}

func (m *MockInterestRepository) FindAll() ([]domain.Interest, error) {
	args := m.Called()
	return args.Get(0).([]domain.Interest), args.Error(1)
}

func (m *MockInterestRepository) FindByID(id uuid.UUID) (*domain.Interest, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Interest), args.Error(1)
}

func (m *MockInterestRepository) FindByIDs(ids []uuid.UUID) ([]domain.Interest, error) {
	args := m.Called(ids)
	return args.Get(0).([]domain.Interest), args.Error(1)
}

func (m *MockInterestRepository) FindByName(name string) (*domain.Interest, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Interest), args.Error(1)
}

func (m *MockInterestRepository) SeedDefaults() error {
	args := m.Called()
	return args.Error(0)
}

// getTestJWTConfig retorna uma configuração JWT para testes.
func getTestJWTConfig() config.JWTConfig {
	return config.JWTConfig{
		SecretKey:          "test-secret-key",
		AccessTokenExpiry:  24,
		RefreshTokenExpiry: 7,
	}
}

// TestAuthService_Register_Success testa o cadastro com sucesso.
func TestAuthService_Register_Success(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	authService := service.NewAuthService(userRepo, interestRepo, getTestJWTConfig())

	req := service.RegisterRequest{
		Name:     "João Silva",
		Email:    "joao@email.com",
		Password: "senha123",
		Age:      30,
		UserType: "VOLUNTEER",
	}

	userRepo.On("ExistsByEmail", req.Email).Return(false, nil)
	userRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	// Act
	result, err := authService.Register(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, req.Name, result.User.Name)
	assert.Equal(t, req.Email, result.User.Email)
	userRepo.AssertExpectations(t)
}

// TestAuthService_Register_EmailAlreadyExists testa cadastro com email duplicado.
func TestAuthService_Register_EmailAlreadyExists(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	authService := service.NewAuthService(userRepo, interestRepo, getTestJWTConfig())

	req := service.RegisterRequest{
		Name:     "João Silva",
		Email:    "joao@email.com",
		Password: "senha123",
		UserType: "VOLUNTEER",
	}

	userRepo.On("ExistsByEmail", req.Email).Return(true, nil)

	// Act
	result, err := authService.Register(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "email já está em uso", err.Error())
	userRepo.AssertExpectations(t)
}

// TestAuthService_Register_WithInterests testa cadastro com interesses.
func TestAuthService_Register_WithInterests(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	authService := service.NewAuthService(userRepo, interestRepo, getTestJWTConfig())

	interestID := uuid.New()
	req := service.RegisterRequest{
		Name:        "Maria Santos",
		Email:       "maria@email.com",
		Password:    "senha456",
		Age:         70,
		UserType:    "ELDERLY",
		InterestIDs: []uuid.UUID{interestID},
	}

	interests := []domain.Interest{
		{ID: interestID, Name: "Música", Icon: "🎵"},
	}

	userRepo.On("ExistsByEmail", req.Email).Return(false, nil)
	interestRepo.On("FindByIDs", req.InterestIDs).Return(interests, nil)
	userRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	// Act
	result, err := authService.Register(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.User.Interests, 1)
	userRepo.AssertExpectations(t)
	interestRepo.AssertExpectations(t)
}

// TestAuthService_Login_Success testa login com sucesso.
func TestAuthService_Login_Success(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	authService := service.NewAuthService(userRepo, interestRepo, getTestJWTConfig())

	password := "senha123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &domain.User{
		ID:           uuid.New(),
		Name:         "João Silva",
		Email:        "joao@email.com",
		PasswordHash: string(hashedPassword),
		UserType:     domain.UserTypeVolunteer,
		IsActive:     true,
	}

	req := service.LoginRequest{
		Email:    "joao@email.com",
		Password: password,
	}

	userRepo.On("FindByEmail", req.Email).Return(user, nil)

	// Act
	result, err := authService.Login(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.AccessToken)
	assert.Equal(t, user.Email, result.User.Email)
	userRepo.AssertExpectations(t)
}

// TestAuthService_Login_InvalidPassword testa login com senha incorreta.
func TestAuthService_Login_InvalidPassword(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	authService := service.NewAuthService(userRepo, interestRepo, getTestJWTConfig())

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)

	user := &domain.User{
		ID:           uuid.New(),
		Email:        "joao@email.com",
		PasswordHash: string(hashedPassword),
		IsActive:     true,
	}

	req := service.LoginRequest{
		Email:    "joao@email.com",
		Password: "senhaerrada",
	}

	userRepo.On("FindByEmail", req.Email).Return(user, nil)

	// Act
	result, err := authService.Login(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "credenciais inválidas", err.Error())
}

// TestAuthService_Login_UserNotFound testa login com usuário inexistente.
func TestAuthService_Login_UserNotFound(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	authService := service.NewAuthService(userRepo, interestRepo, getTestJWTConfig())

	req := service.LoginRequest{
		Email:    "naoexiste@email.com",
		Password: "senha123",
	}

	userRepo.On("FindByEmail", req.Email).Return(nil, errors.New("usuário não encontrado"))

	// Act
	result, err := authService.Login(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "credenciais inválidas", err.Error())
}

// TestAuthService_Login_InactiveUser testa login com usuário desativado.
func TestAuthService_Login_InactiveUser(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	authService := service.NewAuthService(userRepo, interestRepo, getTestJWTConfig())

	password := "senha123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &domain.User{
		ID:           uuid.New(),
		Email:        "joao@email.com",
		PasswordHash: string(hashedPassword),
		IsActive:     false, // Usuário desativado
	}

	req := service.LoginRequest{
		Email:    "joao@email.com",
		Password: password,
	}

	userRepo.On("FindByEmail", req.Email).Return(user, nil)

	// Act
	result, err := authService.Login(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "usuário desativado", err.Error())
}

// TestAuthService_ValidateToken_Success testa validação de token válido.
func TestAuthService_ValidateToken_Success(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	authService := service.NewAuthService(userRepo, interestRepo, getTestJWTConfig())

	// Primeiro registra um usuário para obter um token válido
	req := service.RegisterRequest{
		Name:     "Teste",
		Email:    "teste@email.com",
		Password: "senha123",
		UserType: "VOLUNTEER",
	}

	userRepo.On("ExistsByEmail", req.Email).Return(false, nil)
	userRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	result, _ := authService.Register(req)

	// Act
	claims, err := authService.ValidateToken(result.AccessToken)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, req.Email, claims.Email)
}

// TestAuthService_ValidateToken_Invalid testa validação de token inválido.
func TestAuthService_ValidateToken_Invalid(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	interestRepo := new(MockInterestRepository)
	authService := service.NewAuthService(userRepo, interestRepo, getTestJWTConfig())

	// Act
	claims, err := authService.ValidateToken("token-invalido")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, "token inválido", err.Error())
}
