// Package service_test contém os testes dos serviços da aplicação.
package service_test

import (
	"testing"

	"amigos-terceira-idade/internal/domain"
	"amigos-terceira-idade/internal/repository"
	"amigos-terceira-idade/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockConnectionRepository implementa repository.ConnectionRepositoryInterface para testes.
type MockConnectionRepository struct {
	mock.Mock
}

// Garante que implementa a interface
var _ repository.ConnectionRepositoryInterface = (*MockConnectionRepository)(nil)

func (m *MockConnectionRepository) Create(connection *domain.Connection) error {
	args := m.Called(connection)
	return args.Error(0)
}

func (m *MockConnectionRepository) FindByID(id uuid.UUID) (*domain.Connection, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Connection), args.Error(1)
}

func (m *MockConnectionRepository) FindByVolunteerID(volunteerID uuid.UUID) ([]domain.Connection, error) {
	args := m.Called(volunteerID)
	return args.Get(0).([]domain.Connection), args.Error(1)
}

func (m *MockConnectionRepository) FindByTargetID(targetID uuid.UUID) ([]domain.Connection, error) {
	args := m.Called(targetID)
	return args.Get(0).([]domain.Connection), args.Error(1)
}

func (m *MockConnectionRepository) FindAcceptedByVolunteer(volunteerID uuid.UUID) ([]domain.Connection, error) {
	args := m.Called(volunteerID)
	return args.Get(0).([]domain.Connection), args.Error(1)
}

func (m *MockConnectionRepository) Exists(volunteerID, targetID uuid.UUID) (bool, error) {
	args := m.Called(volunteerID, targetID)
	return args.Bool(0), args.Error(1)
}

func (m *MockConnectionRepository) Update(connection *domain.Connection) error {
	args := m.Called(connection)
	return args.Error(0)
}

func (m *MockConnectionRepository) UpdateStatus(id uuid.UUID, status domain.ConnectionStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockConnectionRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestMatchingService_GetSuggestions_Success testa busca de sugestões com sucesso.
func TestMatchingService_GetSuggestions_Success(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	connectionRepo := new(MockConnectionRepository)
	matchingService := service.NewMatchingService(userRepo, connectionRepo)

	volunteerID := uuid.New()
	volunteer := &domain.User{
		ID:       volunteerID,
		Name:     "Ricardo",
		UserType: domain.UserTypeVolunteer,
		Interests: []domain.Interest{
			{ID: uuid.New(), Name: "Música"},
			{ID: uuid.New(), Name: "Xadrez"},
		},
	}

	elderlyID := uuid.New()
	elderly := []domain.User{
		{
			ID:       elderlyID,
			Name:     "Gerson",
			UserType: domain.UserTypeElderly,
			Interests: []domain.Interest{
				{ID: volunteer.Interests[0].ID, Name: "Música"}, // Match!
			},
		},
	}

	userRepo.On("FindByID", volunteerID).Return(volunteer, nil)
	userRepo.On("FindByType", domain.UserTypeElderly).Return(elderly, nil)
	userRepo.On("FindByType", domain.UserTypeInstitution).Return([]domain.User{}, nil)
	connectionRepo.On("Exists", volunteerID, elderlyID).Return(false, nil)

	// Act
	suggestions, err := matchingService.GetSuggestions(volunteerID, "")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, suggestions, 1)
	assert.Equal(t, 1, suggestions[0].MatchedInterests)
	userRepo.AssertExpectations(t)
}

// TestMatchingService_GetSuggestions_FilterByType testa busca filtrada por tipo.
func TestMatchingService_GetSuggestions_FilterByType(t *testing.T) {
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

	institutionID := uuid.New()
	institutions := []domain.User{
		{ID: institutionID, Name: "Lar São Vicente", UserType: domain.UserTypeInstitution},
	}

	userRepo.On("FindByID", volunteerID).Return(volunteer, nil)
	userRepo.On("FindByType", domain.UserTypeInstitution).Return(institutions, nil)
	connectionRepo.On("Exists", volunteerID, institutionID).Return(false, nil)

	// Act - filtra apenas instituições
	suggestions, err := matchingService.GetSuggestions(volunteerID, "institution")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, suggestions, 1)
	assert.Equal(t, "Lar São Vicente", suggestions[0].User.Name)
}

// TestMatchingService_GetSuggestions_NotVolunteer testa erro quando usuário não é voluntário.
func TestMatchingService_GetSuggestions_NotVolunteer(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	connectionRepo := new(MockConnectionRepository)
	matchingService := service.NewMatchingService(userRepo, connectionRepo)

	elderlyID := uuid.New()
	elderly := &domain.User{
		ID:       elderlyID,
		UserType: domain.UserTypeElderly, // Não é voluntário
	}

	userRepo.On("FindByID", elderlyID).Return(elderly, nil)

	// Act
	suggestions, err := matchingService.GetSuggestions(elderlyID, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, suggestions)
	assert.Equal(t, "apenas voluntários podem buscar conexões", err.Error())
}

// TestMatchingService_Connect_Success testa criação de conexão com sucesso.
func TestMatchingService_Connect_Success(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	connectionRepo := new(MockConnectionRepository)
	matchingService := service.NewMatchingService(userRepo, connectionRepo)

	volunteerID := uuid.New()
	targetID := uuid.New()

	volunteer := &domain.User{
		ID:        volunteerID,
		UserType:  domain.UserTypeVolunteer,
		Interests: []domain.Interest{{ID: uuid.New(), Name: "Música"}},
	}

	target := &domain.User{
		ID:        targetID,
		UserType:  domain.UserTypeElderly,
		Interests: []domain.Interest{{ID: volunteer.Interests[0].ID, Name: "Música"}},
	}

	connectionRepo.On("Exists", volunteerID, targetID).Return(false, nil)
	userRepo.On("FindByID", volunteerID).Return(volunteer, nil)
	userRepo.On("FindByID", targetID).Return(target, nil)
	connectionRepo.On("Create", mock.AnythingOfType("*domain.Connection")).Return(nil)

	// Act
	connection, err := matchingService.Connect(volunteerID, targetID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, connection)
	assert.Equal(t, volunteerID, connection.VolunteerID)
	assert.Equal(t, targetID, connection.TargetID)
	assert.Equal(t, 1, connection.MatchedInterests)
	connectionRepo.AssertExpectations(t)
}

// TestMatchingService_Connect_AlreadyExists testa erro quando conexão já existe.
func TestMatchingService_Connect_AlreadyExists(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	connectionRepo := new(MockConnectionRepository)
	matchingService := service.NewMatchingService(userRepo, connectionRepo)

	volunteerID := uuid.New()
	targetID := uuid.New()

	connectionRepo.On("Exists", volunteerID, targetID).Return(true, nil)

	// Act
	connection, err := matchingService.Connect(volunteerID, targetID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, connection)
	assert.Equal(t, "conexão já existe", err.Error())
}

// TestMatchingService_Connect_CannotConnectWithVolunteer testa erro ao conectar com voluntário.
func TestMatchingService_Connect_CannotConnectWithVolunteer(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	connectionRepo := new(MockConnectionRepository)
	matchingService := service.NewMatchingService(userRepo, connectionRepo)

	volunteerID := uuid.New()
	targetID := uuid.New()

	volunteer := &domain.User{ID: volunteerID, UserType: domain.UserTypeVolunteer}
	target := &domain.User{ID: targetID, UserType: domain.UserTypeVolunteer} // Também é voluntário

	connectionRepo.On("Exists", volunteerID, targetID).Return(false, nil)
	userRepo.On("FindByID", volunteerID).Return(volunteer, nil)
	userRepo.On("FindByID", targetID).Return(target, nil)

	// Act
	connection, err := matchingService.Connect(volunteerID, targetID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, connection)
	assert.Equal(t, "não é possível conectar com outro voluntário", err.Error())
}

// TestMatchingService_AcceptConnection_Success testa aceitar conexão com sucesso.
func TestMatchingService_AcceptConnection_Success(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	connectionRepo := new(MockConnectionRepository)
	matchingService := service.NewMatchingService(userRepo, connectionRepo)

	connectionID := uuid.New()
	connectionRepo.On("UpdateStatus", connectionID, domain.ConnectionStatusAccepted).Return(nil)

	// Act
	err := matchingService.AcceptConnection(connectionID)

	// Assert
	assert.NoError(t, err)
	connectionRepo.AssertExpectations(t)
}

// TestMatchingService_RejectConnection_Success testa rejeitar conexão com sucesso.
func TestMatchingService_RejectConnection_Success(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	connectionRepo := new(MockConnectionRepository)
	matchingService := service.NewMatchingService(userRepo, connectionRepo)

	connectionID := uuid.New()
	connectionRepo.On("UpdateStatus", connectionID, domain.ConnectionStatusRejected).Return(nil)

	// Act
	err := matchingService.RejectConnection(connectionID)

	// Assert
	assert.NoError(t, err)
	connectionRepo.AssertExpectations(t)
}

// TestMatchingService_GetConnections_Volunteer testa listar conexões de voluntário.
func TestMatchingService_GetConnections_Volunteer(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	connectionRepo := new(MockConnectionRepository)
	matchingService := service.NewMatchingService(userRepo, connectionRepo)

	volunteerID := uuid.New()
	volunteer := &domain.User{ID: volunteerID, UserType: domain.UserTypeVolunteer}
	connections := []domain.Connection{
		{ID: uuid.New(), VolunteerID: volunteerID, Status: domain.ConnectionStatusAccepted},
	}

	userRepo.On("FindByID", volunteerID).Return(volunteer, nil)
	connectionRepo.On("FindByVolunteerID", volunteerID).Return(connections, nil)

	// Act
	result, err := matchingService.GetConnections(volunteerID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

// TestMatchingService_GetConnections_Elderly testa listar conexões de idoso.
func TestMatchingService_GetConnections_Elderly(t *testing.T) {
	// Arrange
	userRepo := new(MockUserRepository)
	connectionRepo := new(MockConnectionRepository)
	matchingService := service.NewMatchingService(userRepo, connectionRepo)

	elderlyID := uuid.New()
	elderly := &domain.User{ID: elderlyID, UserType: domain.UserTypeElderly}
	connections := []domain.Connection{
		{ID: uuid.New(), TargetID: elderlyID, Status: domain.ConnectionStatusPending},
	}

	userRepo.On("FindByID", elderlyID).Return(elderly, nil)
	connectionRepo.On("FindByTargetID", elderlyID).Return(connections, nil)

	// Act
	result, err := matchingService.GetConnections(elderlyID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}
