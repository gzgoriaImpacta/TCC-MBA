// Package repository contém as implementações de acesso a dados.
// Este arquivo define as interfaces dos repositórios para permitir mocks em testes.
package repository

import (
	"amigos-terceira-idade/internal/domain"
	"github.com/google/uuid"
)

// UserRepositoryInterface define as operações do repositório de usuários.
// Permite injeção de dependência e mocks para testes.
type UserRepositoryInterface interface {
	Create(user *domain.User) error
	FindByID(id uuid.UUID) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	ExistsByEmail(email string) (bool, error)
	Update(user *domain.User) error
	Delete(id uuid.UUID) error
	FindByType(userType domain.UserType) ([]domain.User, error)
	AddInterests(userID uuid.UUID, interests []domain.Interest) error
	RemoveInterest(userID uuid.UUID, interestID uuid.UUID) error
	UpdateInterests(userID uuid.UUID, interests []domain.Interest) error
}

// InterestRepositoryInterface define as operações do repositório de interesses.
type InterestRepositoryInterface interface {
	Create(interest *domain.Interest) error
	FindAll() ([]domain.Interest, error)
	FindByID(id uuid.UUID) (*domain.Interest, error)
	FindByIDs(ids []uuid.UUID) ([]domain.Interest, error)
	FindByName(name string) (*domain.Interest, error)
	SeedDefaults() error
}

// ConnectionRepositoryInterface define as operações do repositório de conexões.
type ConnectionRepositoryInterface interface {
	Create(connection *domain.Connection) error
	FindByID(id uuid.UUID) (*domain.Connection, error)
	FindByVolunteerID(volunteerID uuid.UUID) ([]domain.Connection, error)
	FindByTargetID(targetID uuid.UUID) ([]domain.Connection, error)
	FindAcceptedByVolunteer(volunteerID uuid.UUID) ([]domain.Connection, error)
	Exists(volunteerID, targetID uuid.UUID) (bool, error)
	Update(connection *domain.Connection) error
	UpdateStatus(id uuid.UUID, status domain.ConnectionStatus) error
	Delete(id uuid.UUID) error
}

// AppointmentRepositoryInterface define as operações do repositório de agendamentos.
type AppointmentRepositoryInterface interface {
	Create(appointment *domain.Appointment) error
	FindByID(id uuid.UUID) (*domain.Appointment, error)
	FindByVolunteerID(volunteerID uuid.UUID) ([]domain.Appointment, error)
	FindByTargetID(targetID uuid.UUID) ([]domain.Appointment, error)
	FindUpcoming(userID uuid.UUID) ([]domain.Appointment, error)
	FindPendingInvitations(targetID uuid.UUID) ([]domain.Appointment, error)
	FindSentInvitations(volunteerID uuid.UUID) ([]domain.Appointment, error)
	Update(appointment *domain.Appointment) error
	UpdateStatus(id uuid.UUID, status domain.AppointmentStatus) error
	Delete(id uuid.UUID) error
}

// Garante que as implementações satisfazem as interfaces
var _ UserRepositoryInterface = (*UserRepository)(nil)
var _ InterestRepositoryInterface = (*InterestRepository)(nil)
var _ ConnectionRepositoryInterface = (*ConnectionRepository)(nil)
var _ AppointmentRepositoryInterface = (*AppointmentRepository)(nil)
