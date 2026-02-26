// Package domain contém as entidades de negócio da aplicação.
package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConnectionStatus define os possíveis estados de uma conexão.
type ConnectionStatus string

const (
	ConnectionStatusPending  ConnectionStatus = "PENDING"  // Aguardando aceite
	ConnectionStatusAccepted ConnectionStatus = "ACCEPTED" // Conexão aceita
	ConnectionStatusRejected ConnectionStatus = "REJECTED" // Conexão rejeitada
)

// Connection representa uma conexão/pareamento entre voluntário e idoso/instituição.
type Connection struct {
	ID               uuid.UUID        `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	VolunteerID      uuid.UUID        `gorm:"type:uniqueidentifier;not null" json:"volunteer_id"`
	TargetID         uuid.UUID        `gorm:"type:uniqueidentifier;not null" json:"target_id"`
	TargetType       UserType         `gorm:"size:20;not null" json:"target_type"` // ELDERLY ou INSTITUTION
	Status           ConnectionStatus `gorm:"size:20;default:PENDING" json:"status"`
	MatchedInterests int              `gorm:"default:0" json:"matched_interests"` // Quantidade de interesses em comum
	CreatedAt        time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time        `gorm:"autoUpdateTime" json:"updated_at"`

	// Relacionamentos para facilitar consultas
	Volunteer User `gorm:"foreignKey:VolunteerID" json:"volunteer,omitempty"`
	Target    User `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}

// TableName define o nome da tabela no banco de dados.
func (Connection) TableName() string {
	return "connections"
}

// BeforeCreate é executado antes de inserir uma nova conexão.
func (c *Connection) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
