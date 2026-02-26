// Package service_test contém os testes dos serviços da aplicação.
package service_test

import (
	"testing"
	"time"

	"amigos-terceira-idade/internal/domain"
	"amigos-terceira-idade/internal/repository"
	"amigos-terceira-idade/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAppointmentRepository implementa repository.AppointmentRepositoryInterface para testes.
type MockAppointmentRepository struct {
	mock.Mock
}

// Garante que implementa a interface
var _ repository.AppointmentRepositoryInterface = (*MockAppointmentRepository)(nil)

func (m *MockAppointmentRepository) Create(appointment *domain.Appointment) error {
	args := m.Called(appointment)
	return args.Error(0)
}

func (m *MockAppointmentRepository) FindByID(id uuid.UUID) (*domain.Appointment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) FindByVolunteerID(volunteerID uuid.UUID) ([]domain.Appointment, error) {
	args := m.Called(volunteerID)
	return args.Get(0).([]domain.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) FindByTargetID(targetID uuid.UUID) ([]domain.Appointment, error) {
	args := m.Called(targetID)
	return args.Get(0).([]domain.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) FindUpcoming(userID uuid.UUID) ([]domain.Appointment, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) FindPendingInvitations(targetID uuid.UUID) ([]domain.Appointment, error) {
	args := m.Called(targetID)
	return args.Get(0).([]domain.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) FindSentInvitations(volunteerID uuid.UUID) ([]domain.Appointment, error) {
	args := m.Called(volunteerID)
	return args.Get(0).([]domain.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) Update(appointment *domain.Appointment) error {
	args := m.Called(appointment)
	return args.Error(0)
}

func (m *MockAppointmentRepository) UpdateStatus(id uuid.UUID, status domain.AppointmentStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockAppointmentRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestAppointmentService_Create_Success testa criação de agendamento com sucesso.
func TestAppointmentService_Create_Success(t *testing.T) {
	// Arrange
	appointmentRepo := new(MockAppointmentRepository)
	userRepo := new(MockUserRepository)
	appointmentService := service.NewAppointmentService(appointmentRepo, userRepo)

	volunteerID := uuid.New()
	targetID := uuid.New()

	volunteer := &domain.User{ID: volunteerID, UserType: domain.UserTypeVolunteer}
	target := &domain.User{ID: targetID, UserType: domain.UserTypeElderly}

	futureDate := time.Now().Add(24 * time.Hour)
	req := service.CreateAppointmentRequest{
		TargetID:        targetID,
		Date:            futureDate,
		DurationMinutes: 30,
	}

	userRepo.On("FindByID", volunteerID).Return(volunteer, nil)
	userRepo.On("FindByID", targetID).Return(target, nil)
	appointmentRepo.On("Create", mock.AnythingOfType("*domain.Appointment")).Return(nil)
	appointmentRepo.On("FindByID", mock.AnythingOfType("uuid.UUID")).Return(&domain.Appointment{
		ID:              uuid.New(),
		VolunteerID:     volunteerID,
		TargetID:        targetID,
		Date:            futureDate,
		DurationMinutes: 30,
		Status:          domain.AppointmentStatusPending,
	}, nil)

	// Act
	result, err := appointmentService.Create(volunteerID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, targetID, result.TargetID)
	assert.Equal(t, domain.AppointmentStatusPending, result.Status)
}

// TestAppointmentService_Create_NotVolunteer testa erro quando criador não é voluntário.
func TestAppointmentService_Create_NotVolunteer(t *testing.T) {
	// Arrange
	appointmentRepo := new(MockAppointmentRepository)
	userRepo := new(MockUserRepository)
	appointmentService := service.NewAppointmentService(appointmentRepo, userRepo)

	elderlyID := uuid.New()
	elderly := &domain.User{ID: elderlyID, UserType: domain.UserTypeElderly}

	req := service.CreateAppointmentRequest{
		TargetID: uuid.New(),
		Date:     time.Now().Add(24 * time.Hour),
	}

	userRepo.On("FindByID", elderlyID).Return(elderly, nil)

	// Act
	result, err := appointmentService.Create(elderlyID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "apenas voluntários podem criar agendamentos", err.Error())
}

// TestAppointmentService_Create_PastDate testa erro quando data é passada.
func TestAppointmentService_Create_PastDate(t *testing.T) {
	// Arrange
	appointmentRepo := new(MockAppointmentRepository)
	userRepo := new(MockUserRepository)
	appointmentService := service.NewAppointmentService(appointmentRepo, userRepo)

	volunteerID := uuid.New()
	targetID := uuid.New()

	volunteer := &domain.User{ID: volunteerID, UserType: domain.UserTypeVolunteer}
	target := &domain.User{ID: targetID, UserType: domain.UserTypeElderly}

	pastDate := time.Now().Add(-24 * time.Hour) // Data passada
	req := service.CreateAppointmentRequest{
		TargetID: targetID,
		Date:     pastDate,
	}

	userRepo.On("FindByID", volunteerID).Return(volunteer, nil)
	userRepo.On("FindByID", targetID).Return(target, nil)

	// Act
	result, err := appointmentService.Create(volunteerID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "a data deve ser futura", err.Error())
}

// TestAppointmentService_Accept_Success testa aceitar convite com sucesso.
func TestAppointmentService_Accept_Success(t *testing.T) {
	// Arrange
	appointmentRepo := new(MockAppointmentRepository)
	userRepo := new(MockUserRepository)
	appointmentService := service.NewAppointmentService(appointmentRepo, userRepo)

	appointmentID := uuid.New()
	targetID := uuid.New()

	appointment := &domain.Appointment{
		ID:       appointmentID,
		TargetID: targetID,
		Status:   domain.AppointmentStatusPending,
	}

	appointmentRepo.On("FindByID", appointmentID).Return(appointment, nil)
	appointmentRepo.On("UpdateStatus", appointmentID, domain.AppointmentStatusConfirmed).Return(nil)

	// Act
	err := appointmentService.Accept(appointmentID, targetID)

	// Assert
	assert.NoError(t, err)
	appointmentRepo.AssertExpectations(t)
}

// TestAppointmentService_Accept_NotTarget testa erro quando usuário não é destinatário.
func TestAppointmentService_Accept_NotTarget(t *testing.T) {
	// Arrange
	appointmentRepo := new(MockAppointmentRepository)
	userRepo := new(MockUserRepository)
	appointmentService := service.NewAppointmentService(appointmentRepo, userRepo)

	appointmentID := uuid.New()
	targetID := uuid.New()
	wrongUserID := uuid.New()

	appointment := &domain.Appointment{
		ID:       appointmentID,
		TargetID: targetID,
		Status:   domain.AppointmentStatusPending,
	}

	appointmentRepo.On("FindByID", appointmentID).Return(appointment, nil)

	// Act
	err := appointmentService.Accept(appointmentID, wrongUserID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "você não pode aceitar este convite", err.Error())
}

// TestAppointmentService_Accept_NotPending testa erro quando convite não está pendente.
func TestAppointmentService_Accept_NotPending(t *testing.T) {
	// Arrange
	appointmentRepo := new(MockAppointmentRepository)
	userRepo := new(MockUserRepository)
	appointmentService := service.NewAppointmentService(appointmentRepo, userRepo)

	appointmentID := uuid.New()
	targetID := uuid.New()

	appointment := &domain.Appointment{
		ID:       appointmentID,
		TargetID: targetID,
		Status:   domain.AppointmentStatusConfirmed, // Já confirmado
	}

	appointmentRepo.On("FindByID", appointmentID).Return(appointment, nil)

	// Act
	err := appointmentService.Accept(appointmentID, targetID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "este convite não está mais pendente", err.Error())
}

// TestAppointmentService_Decline_Success testa recusar convite com sucesso.
func TestAppointmentService_Decline_Success(t *testing.T) {
	// Arrange
	appointmentRepo := new(MockAppointmentRepository)
	userRepo := new(MockUserRepository)
	appointmentService := service.NewAppointmentService(appointmentRepo, userRepo)

	appointmentID := uuid.New()
	targetID := uuid.New()

	appointment := &domain.Appointment{
		ID:       appointmentID,
		TargetID: targetID,
		Status:   domain.AppointmentStatusPending,
	}

	appointmentRepo.On("FindByID", appointmentID).Return(appointment, nil)
	appointmentRepo.On("UpdateStatus", appointmentID, domain.AppointmentStatusCancelled).Return(nil)

	// Act
	err := appointmentService.Decline(appointmentID, targetID)

	// Assert
	assert.NoError(t, err)
	appointmentRepo.AssertExpectations(t)
}

// TestAppointmentService_Cancel_Success testa cancelar agendamento com sucesso.
func TestAppointmentService_Cancel_Success(t *testing.T) {
	// Arrange
	appointmentRepo := new(MockAppointmentRepository)
	userRepo := new(MockUserRepository)
	appointmentService := service.NewAppointmentService(appointmentRepo, userRepo)

	appointmentID := uuid.New()
	volunteerID := uuid.New()

	appointment := &domain.Appointment{
		ID:          appointmentID,
		VolunteerID: volunteerID,
		Status:      domain.AppointmentStatusConfirmed,
	}

	appointmentRepo.On("FindByID", appointmentID).Return(appointment, nil)
	appointmentRepo.On("UpdateStatus", appointmentID, domain.AppointmentStatusCancelled).Return(nil)

	// Act
	err := appointmentService.Cancel(appointmentID, volunteerID)

	// Assert
	assert.NoError(t, err)
}

// TestAppointmentService_Cancel_NotParticipant testa erro quando usuário não é participante.
func TestAppointmentService_Cancel_NotParticipant(t *testing.T) {
	// Arrange
	appointmentRepo := new(MockAppointmentRepository)
	userRepo := new(MockUserRepository)
	appointmentService := service.NewAppointmentService(appointmentRepo, userRepo)

	appointmentID := uuid.New()
	volunteerID := uuid.New()
	targetID := uuid.New()
	wrongUserID := uuid.New()

	appointment := &domain.Appointment{
		ID:          appointmentID,
		VolunteerID: volunteerID,
		TargetID:    targetID,
	}

	appointmentRepo.On("FindByID", appointmentID).Return(appointment, nil)

	// Act
	err := appointmentService.Cancel(appointmentID, wrongUserID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "você não pode cancelar este agendamento", err.Error())
}

// TestAppointmentService_GetMyAppointments_Volunteer testa listar agendamentos de voluntário.
func TestAppointmentService_GetMyAppointments_Volunteer(t *testing.T) {
	// Arrange
	appointmentRepo := new(MockAppointmentRepository)
	userRepo := new(MockUserRepository)
	appointmentService := service.NewAppointmentService(appointmentRepo, userRepo)

	volunteerID := uuid.New()
	volunteer := &domain.User{ID: volunteerID, UserType: domain.UserTypeVolunteer}
	appointments := []domain.Appointment{
		{ID: uuid.New(), VolunteerID: volunteerID},
		{ID: uuid.New(), VolunteerID: volunteerID},
	}

	userRepo.On("FindByID", volunteerID).Return(volunteer, nil)
	appointmentRepo.On("FindByVolunteerID", volunteerID).Return(appointments, nil)

	// Act
	result, err := appointmentService.GetMyAppointments(volunteerID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

// TestAppointmentService_GetUpcoming_Success testa listar próximos agendamentos.
func TestAppointmentService_GetUpcoming_Success(t *testing.T) {
	// Arrange
	appointmentRepo := new(MockAppointmentRepository)
	userRepo := new(MockUserRepository)
	appointmentService := service.NewAppointmentService(appointmentRepo, userRepo)

	userID := uuid.New()
	futureDate := time.Now().Add(24 * time.Hour)
	appointments := []domain.Appointment{
		{ID: uuid.New(), Date: futureDate, Status: domain.AppointmentStatusConfirmed},
	}

	appointmentRepo.On("FindUpcoming", userID).Return(appointments, nil)

	// Act
	result, err := appointmentService.GetUpcoming(userID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

// TestAppointmentService_Complete_WithRating testa marcar como completo com avaliação.
func TestAppointmentService_Complete_WithRating(t *testing.T) {
	// Arrange
	appointmentRepo := new(MockAppointmentRepository)
	userRepo := new(MockUserRepository)
	appointmentService := service.NewAppointmentService(appointmentRepo, userRepo)

	appointmentID := uuid.New()
	appointment := &domain.Appointment{
		ID:     appointmentID,
		Status: domain.AppointmentStatusConfirmed,
	}

	appointmentRepo.On("FindByID", appointmentID).Return(appointment, nil)
	appointmentRepo.On("Update", mock.AnythingOfType("*domain.Appointment")).Return(nil)

	// Act
	err := appointmentService.Complete(appointmentID, 5)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, domain.AppointmentStatusCompleted, appointment.Status)
	assert.Equal(t, 5, appointment.Rating)
}
