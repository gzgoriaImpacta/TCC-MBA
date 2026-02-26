// Package service_test contém os testes adicionais de appointment.
package service_test

import (
	"testing"

	"amigos-terceira-idade/internal/domain"
	"amigos-terceira-idade/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestAppointmentService_SetMeetingURL_Success testa definir URL de reunião.
func TestAppointmentService_SetMeetingURL_Success(t *testing.T) {
	// Arrange
	appointmentRepo := new(MockAppointmentRepository)
	userRepo := new(MockUserRepository)
	appointmentService := service.NewAppointmentService(appointmentRepo, userRepo)

	appointmentID := uuid.New()
	appointment := &domain.Appointment{
		ID:     appointmentID,
		Status: domain.AppointmentStatusConfirmed,
	}
	meetingURL := "https://meet.google.com/abc-defg-hij"

	appointmentRepo.On("FindByID", appointmentID).Return(appointment, nil)
	appointmentRepo.On("Update", mock.AnythingOfType("*domain.Appointment")).Return(nil)

	// Act
	err := appointmentService.SetMeetingURL(appointmentID, meetingURL)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, meetingURL, appointment.MeetingURL)
}

// TestAppointmentService_GetReceivedInvitations_Success testa listar convites recebidos.
func TestAppointmentService_GetReceivedInvitations_Success(t *testing.T) {
	// Arrange
	appointmentRepo := new(MockAppointmentRepository)
	userRepo := new(MockUserRepository)
	appointmentService := service.NewAppointmentService(appointmentRepo, userRepo)

	userID := uuid.New()
	invitations := []domain.Appointment{
		{ID: uuid.New(), TargetID: userID, Status: domain.AppointmentStatusPending},
	}

	appointmentRepo.On("FindPendingInvitations", userID).Return(invitations, nil)

	// Act
	result, err := appointmentService.GetReceivedInvitations(userID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

// TestAppointmentService_GetSentInvitations_Success testa listar convites enviados.
func TestAppointmentService_GetSentInvitations_Success(t *testing.T) {
	// Arrange
	appointmentRepo := new(MockAppointmentRepository)
	userRepo := new(MockUserRepository)
	appointmentService := service.NewAppointmentService(appointmentRepo, userRepo)

	volunteerID := uuid.New()
	invitations := []domain.Appointment{
		{ID: uuid.New(), VolunteerID: volunteerID, Status: domain.AppointmentStatusPending},
	}

	appointmentRepo.On("FindSentInvitations", volunteerID).Return(invitations, nil)

	// Act
	result, err := appointmentService.GetSentInvitations(volunteerID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

// TestAppointmentService_GetByID_Success testa buscar agendamento por ID.
func TestAppointmentService_GetByID_Success(t *testing.T) {
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

	// Act
	result, err := appointmentService.GetByID(appointmentID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, appointmentID, result.ID)
}

// TestAppointmentService_Decline_NotTarget testa recusar convite por não ser destinatário.
func TestAppointmentService_Decline_NotTarget(t *testing.T) {
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
	err := appointmentService.Decline(appointmentID, wrongUserID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "você não pode recusar este convite", err.Error())
}

// TestAppointmentService_GetMyAppointments_Elderly testa listar agendamentos de idoso.
func TestAppointmentService_GetMyAppointments_Elderly(t *testing.T) {
	// Arrange
	appointmentRepo := new(MockAppointmentRepository)
	userRepo := new(MockUserRepository)
	appointmentService := service.NewAppointmentService(appointmentRepo, userRepo)

	elderlyID := uuid.New()
	elderly := &domain.User{ID: elderlyID, UserType: domain.UserTypeElderly}
	appointments := []domain.Appointment{
		{ID: uuid.New(), TargetID: elderlyID},
	}

	userRepo.On("FindByID", elderlyID).Return(elderly, nil)
	appointmentRepo.On("FindByTargetID", elderlyID).Return(appointments, nil)

	// Act
	result, err := appointmentService.GetMyAppointments(elderlyID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}
