// Package service contém a lógica de negócio da aplicação.
package service

import (
	"errors"
	"time"

	"amigos-terceira-idade/internal/domain"
	"amigos-terceira-idade/internal/repository"
	"github.com/google/uuid"
)

// AppointmentService gerencia os agendamentos de conversas.
type AppointmentService struct {
	appointmentRepo repository.AppointmentRepositoryInterface
	userRepo        repository.UserRepositoryInterface
}

// NewAppointmentService cria uma nova instância do serviço de agendamentos.
func NewAppointmentService(
	appointmentRepo repository.AppointmentRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
) *AppointmentService {
	return &AppointmentService{
		appointmentRepo: appointmentRepo,
		userRepo:        userRepo,
	}
}

// CreateAppointmentRequest contém os dados para criar um agendamento.
type CreateAppointmentRequest struct {
	TargetID        uuid.UUID `json:"target_id" binding:"required"`
	Date            time.Time `json:"date" binding:"required"`
	DurationMinutes int       `json:"duration_minutes"`
	Notes           string    `json:"notes"`
}

// Create cria um novo agendamento (envia convite).
func (s *AppointmentService) Create(volunteerID uuid.UUID, req CreateAppointmentRequest) (*domain.Appointment, error) {
	// Valida o voluntário
	volunteer, err := s.userRepo.FindByID(volunteerID)
	if err != nil {
		return nil, err
	}
	if volunteer.UserType != domain.UserTypeVolunteer {
		return nil, errors.New("apenas voluntários podem criar agendamentos")
	}

	// Valida o destinatário
	target, err := s.userRepo.FindByID(req.TargetID)
	if err != nil {
		return nil, err
	}
	if target.UserType == domain.UserTypeVolunteer {
		return nil, errors.New("não é possível agendar com outro voluntário")
	}

	// Valida a data (deve ser futura)
	if req.Date.Before(time.Now()) {
		return nil, errors.New("a data deve ser futura")
	}

	// Define duração padrão se não informada
	duration := req.DurationMinutes
	if duration <= 0 {
		duration = 30 // 30 minutos padrão
	}

	// Cria o agendamento
	appointment := &domain.Appointment{
		ID:              uuid.New(),
		VolunteerID:     volunteerID,
		TargetID:        req.TargetID,
		TargetType:      target.UserType,
		Date:            req.Date,
		DurationMinutes: duration,
		Status:          domain.AppointmentStatusPending,
		Notes:           req.Notes,
	}

	if err := s.appointmentRepo.Create(appointment); err != nil {
		return nil, err
	}

	// Retorna com os relacionamentos preenchidos
	return s.appointmentRepo.FindByID(appointment.ID)
}

// GetByID busca um agendamento pelo ID.
func (s *AppointmentService) GetByID(id uuid.UUID) (*domain.Appointment, error) {
	return s.appointmentRepo.FindByID(id)
}

// GetMyAppointments retorna os agendamentos de um usuário.
func (s *AppointmentService) GetMyAppointments(userID uuid.UUID) ([]domain.Appointment, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user.UserType == domain.UserTypeVolunteer {
		return s.appointmentRepo.FindByVolunteerID(userID)
	}
	return s.appointmentRepo.FindByTargetID(userID)
}

// GetUpcoming retorna os próximos agendamentos confirmados.
func (s *AppointmentService) GetUpcoming(userID uuid.UUID) ([]domain.Appointment, error) {
	return s.appointmentRepo.FindUpcoming(userID)
}

// GetReceivedInvitations retorna os convites recebidos pendentes.
func (s *AppointmentService) GetReceivedInvitations(userID uuid.UUID) ([]domain.Appointment, error) {
	return s.appointmentRepo.FindPendingInvitations(userID)
}

// GetSentInvitations retorna os convites enviados pendentes.
func (s *AppointmentService) GetSentInvitations(userID uuid.UUID) ([]domain.Appointment, error) {
	return s.appointmentRepo.FindSentInvitations(userID)
}

// Accept aceita um convite de agendamento.
func (s *AppointmentService) Accept(appointmentID uuid.UUID, userID uuid.UUID) error {
	appointment, err := s.appointmentRepo.FindByID(appointmentID)
	if err != nil {
		return err
	}

	// Verifica se o usuário é o destinatário do convite
	if appointment.TargetID != userID {
		return errors.New("você não pode aceitar este convite")
	}

	// Verifica se está pendente
	if appointment.Status != domain.AppointmentStatusPending {
		return errors.New("este convite não está mais pendente")
	}

	return s.appointmentRepo.UpdateStatus(appointmentID, domain.AppointmentStatusConfirmed)
}

// Decline recusa um convite de agendamento.
func (s *AppointmentService) Decline(appointmentID uuid.UUID, userID uuid.UUID) error {
	appointment, err := s.appointmentRepo.FindByID(appointmentID)
	if err != nil {
		return err
	}

	// Verifica se o usuário é o destinatário do convite
	if appointment.TargetID != userID {
		return errors.New("você não pode recusar este convite")
	}

	return s.appointmentRepo.UpdateStatus(appointmentID, domain.AppointmentStatusCancelled)
}

// Cancel cancela um agendamento.
func (s *AppointmentService) Cancel(appointmentID uuid.UUID, userID uuid.UUID) error {
	appointment, err := s.appointmentRepo.FindByID(appointmentID)
	if err != nil {
		return err
	}

	// Verifica se o usuário é participante do agendamento
	if appointment.VolunteerID != userID && appointment.TargetID != userID {
		return errors.New("você não pode cancelar este agendamento")
	}

	return s.appointmentRepo.UpdateStatus(appointmentID, domain.AppointmentStatusCancelled)
}

// Complete marca um agendamento como concluído.
func (s *AppointmentService) Complete(appointmentID uuid.UUID, rating int) error {
	appointment, err := s.appointmentRepo.FindByID(appointmentID)
	if err != nil {
		return err
	}

	appointment.Status = domain.AppointmentStatusCompleted
	if rating >= 1 && rating <= 5 {
		appointment.Rating = rating
	}

	return s.appointmentRepo.Update(appointment)
}

// SetMeetingURL define o link da reunião para um agendamento.
func (s *AppointmentService) SetMeetingURL(appointmentID uuid.UUID, url string) error {
	appointment, err := s.appointmentRepo.FindByID(appointmentID)
	if err != nil {
		return err
	}
	appointment.MeetingURL = url
	return s.appointmentRepo.Update(appointment)
}
