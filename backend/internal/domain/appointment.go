// Package domain contém as entidades de negócio da aplicação.
package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AppointmentStatus define os possíveis estados de um agendamento.
type AppointmentStatus string

const (
	AppointmentStatusPending   AppointmentStatus = "PENDING"   // Convite enviado, aguardando aceite
	AppointmentStatusConfirmed AppointmentStatus = "CONFIRMED" // Agendamento confirmado
	AppointmentStatusCancelled AppointmentStatus = "CANCELLED" // Agendamento cancelado
	AppointmentStatusCompleted AppointmentStatus = "COMPLETED" // Conversa realizada
)

// Appointment representa um agendamento de conversa entre voluntário e idoso.
type Appointment struct {
	ID              uuid.UUID         `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	VolunteerID     uuid.UUID         `gorm:"type:uniqueidentifier;not null" json:"volunteer_id"`
	TargetID        uuid.UUID         `gorm:"type:uniqueidentifier;not null" json:"target_id"`
	TargetType      UserType          `gorm:"size:20;not null" json:"target_type"`
	Date            time.Time         `gorm:"not null" json:"date"`
	DurationMinutes int               `gorm:"default:30" json:"duration_minutes"`
	Status          AppointmentStatus `gorm:"size:20;default:PENDING" json:"status"`
	MeetingURL      string            `gorm:"size:500" json:"meeting_url,omitempty"` // Link do Google Meet
	Notes           string            `gorm:"type:text" json:"notes,omitempty"`
	Rating          int               `gorm:"" json:"rating,omitempty"` // Avaliação pós-conversa (1-5)
	CreatedAt       time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time         `gorm:"autoUpdateTime" json:"updated_at"`

	// Relacionamentos
	Volunteer User `gorm:"foreignKey:VolunteerID" json:"volunteer,omitempty"`
	Target    User `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}

// TableName define o nome da tabela no banco de dados.
func (Appointment) TableName() string {
	return "appointments"
}

// BeforeCreate é executado antes de inserir um novo agendamento.
func (a *Appointment) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
