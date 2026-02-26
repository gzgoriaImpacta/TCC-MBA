// Package repository contém as implementações de acesso a dados.
package repository

import (
	"errors"
	"time"

	"amigos-terceira-idade/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AppointmentRepository gerencia as operações de banco de dados para agendamentos.
type AppointmentRepository struct {
	db *gorm.DB
}

// NewAppointmentRepository cria uma nova instância do repositório de agendamentos.
func NewAppointmentRepository(db *gorm.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

// Create insere um novo agendamento no banco de dados.
func (r *AppointmentRepository) Create(appointment *domain.Appointment) error {
	return r.db.Create(appointment).Error
}

// FindByID busca um agendamento pelo ID.
func (r *AppointmentRepository) FindByID(id uuid.UUID) (*domain.Appointment, error) {
	var appointment domain.Appointment
	err := r.db.Preload("Volunteer").Preload("Target").
		First(&appointment, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("agendamento não encontrado")
		}
		return nil, err
	}
	return &appointment, nil
}

// FindByVolunteerID busca todos os agendamentos de um voluntário.
func (r *AppointmentRepository) FindByVolunteerID(volunteerID uuid.UUID) ([]domain.Appointment, error) {
	var appointments []domain.Appointment
	err := r.db.Preload("Target").
		Where("volunteer_id = ?", volunteerID).
		Order("date ASC").
		Find(&appointments).Error
	if err != nil {
		return nil, err
	}
	return appointments, nil
}

// FindByTargetID busca todos os agendamentos de um idoso/instituição.
func (r *AppointmentRepository) FindByTargetID(targetID uuid.UUID) ([]domain.Appointment, error) {
	var appointments []domain.Appointment
	err := r.db.Preload("Volunteer").
		Where("target_id = ?", targetID).
		Order("date ASC").
		Find(&appointments).Error
	if err != nil {
		return nil, err
	}
	return appointments, nil
}

// FindUpcoming busca os próximos agendamentos de um usuário.
// Retorna agendamentos confirmados com data futura.
func (r *AppointmentRepository) FindUpcoming(userID uuid.UUID) ([]domain.Appointment, error) {
	var appointments []domain.Appointment
	now := time.Now()
	err := r.db.Preload("Volunteer").Preload("Target").
		Where("(volunteer_id = ? OR target_id = ?) AND date > ? AND status = ?",
			userID, userID, now, domain.AppointmentStatusConfirmed).
		Order("date ASC").
		Find(&appointments).Error
	if err != nil {
		return nil, err
	}
	return appointments, nil
}

// FindPendingInvitations busca convites pendentes recebidos por um usuário.
func (r *AppointmentRepository) FindPendingInvitations(targetID uuid.UUID) ([]domain.Appointment, error) {
	var appointments []domain.Appointment
	err := r.db.Preload("Volunteer").
		Where("target_id = ? AND status = ?", targetID, domain.AppointmentStatusPending).
		Order("date ASC").
		Find(&appointments).Error
	if err != nil {
		return nil, err
	}
	return appointments, nil
}

// FindSentInvitations busca convites enviados por um voluntário.
func (r *AppointmentRepository) FindSentInvitations(volunteerID uuid.UUID) ([]domain.Appointment, error) {
	var appointments []domain.Appointment
	err := r.db.Preload("Target").
		Where("volunteer_id = ? AND status = ?", volunteerID, domain.AppointmentStatusPending).
		Order("date ASC").
		Find(&appointments).Error
	if err != nil {
		return nil, err
	}
	return appointments, nil
}

// Update atualiza os dados de um agendamento.
func (r *AppointmentRepository) Update(appointment *domain.Appointment) error {
	return r.db.Save(appointment).Error
}

// UpdateStatus atualiza apenas o status de um agendamento.
func (r *AppointmentRepository) UpdateStatus(id uuid.UUID, status domain.AppointmentStatus) error {
	return r.db.Model(&domain.Appointment{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete remove um agendamento.
func (r *AppointmentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Appointment{}, "id = ?", id).Error
}
