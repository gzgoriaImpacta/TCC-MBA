// Package repository contém as implementações de acesso a dados.
package repository

import (
	"errors"

	"amigos-terceira-idade/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConnectionRepository gerencia as operações de banco de dados para conexões/pareamentos.
type ConnectionRepository struct {
	db *gorm.DB
}

// NewConnectionRepository cria uma nova instância do repositório de conexões.
func NewConnectionRepository(db *gorm.DB) *ConnectionRepository {
	return &ConnectionRepository{db: db}
}

// Create insere uma nova conexão no banco de dados.
func (r *ConnectionRepository) Create(connection *domain.Connection) error {
	return r.db.Create(connection).Error
}

// FindByID busca uma conexão pelo ID.
func (r *ConnectionRepository) FindByID(id uuid.UUID) (*domain.Connection, error) {
	var connection domain.Connection
	err := r.db.Preload("Volunteer").Preload("Target").
		First(&connection, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("conexão não encontrada")
		}
		return nil, err
	}
	return &connection, nil
}

// FindByVolunteerID busca todas as conexões de um voluntário.
func (r *ConnectionRepository) FindByVolunteerID(volunteerID uuid.UUID) ([]domain.Connection, error) {
	var connections []domain.Connection
	err := r.db.Preload("Target").
		Where("volunteer_id = ?", volunteerID).
		Find(&connections).Error
	if err != nil {
		return nil, err
	}
	return connections, nil
}

// FindByTargetID busca todas as conexões de um idoso/instituição.
func (r *ConnectionRepository) FindByTargetID(targetID uuid.UUID) ([]domain.Connection, error) {
	var connections []domain.Connection
	err := r.db.Preload("Volunteer").
		Where("target_id = ?", targetID).
		Find(&connections).Error
	if err != nil {
		return nil, err
	}
	return connections, nil
}

// FindAcceptedByVolunteer busca conexões aceitas de um voluntário.
func (r *ConnectionRepository) FindAcceptedByVolunteer(volunteerID uuid.UUID) ([]domain.Connection, error) {
	var connections []domain.Connection
	err := r.db.Preload("Target").
		Where("volunteer_id = ? AND status = ?", volunteerID, domain.ConnectionStatusAccepted).
		Find(&connections).Error
	if err != nil {
		return nil, err
	}
	return connections, nil
}

// Exists verifica se já existe uma conexão entre o voluntário e o alvo.
func (r *ConnectionRepository) Exists(volunteerID, targetID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Connection{}).
		Where("volunteer_id = ? AND target_id = ?", volunteerID, targetID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Update atualiza os dados de uma conexão.
func (r *ConnectionRepository) Update(connection *domain.Connection) error {
	return r.db.Save(connection).Error
}

// UpdateStatus atualiza apenas o status de uma conexão.
func (r *ConnectionRepository) UpdateStatus(id uuid.UUID, status domain.ConnectionStatus) error {
	return r.db.Model(&domain.Connection{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete remove uma conexão.
func (r *ConnectionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Connection{}, "id = ?", id).Error
}
