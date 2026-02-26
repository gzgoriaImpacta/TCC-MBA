// Package domain contém as entidades de negócio da aplicação.
// Estas estruturas representam os objetos principais do sistema.
package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserType define os tipos de usuário no sistema.
type UserType string

const (
	UserTypeVolunteer   UserType = "VOLUNTEER"   // Voluntário que oferece companhia
	UserTypeElderly     UserType = "ELDERLY"     // Idoso que recebe companhia
	UserTypeInstitution UserType = "INSTITUTION" // Instituição/lar de idosos
)

// User representa um usuário do sistema.
// É a entidade base para voluntários, idosos e instituições.
type User struct {
	ID           uuid.UUID `gorm:"type:uniqueidentifier;default:NEWID();primaryKey"`
	Name         string    `gorm:"size:255;not null" json:"name"`
	Email        string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"` // Nunca retorna no JSON
	Age          int       `gorm:"" json:"age,omitempty"`
	Bio          string    `gorm:"type:text" json:"bio,omitempty"`
	Phone        string    `gorm:"type:text" json:"phone,omitempty"`
	PhotoURL     string    `gorm:"size:500" json:"photo_url,omitempty"`
	UserType     UserType  `gorm:"size:20;not null" json:"user_type"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relacionamentos
	Interests []Interest `gorm:"many2many:user_interests;" json:"interests,omitempty"`
}

// TableName define o nome da tabela no banco de dados.
func (User) TableName() string {
	return "users"
}

// BeforeCreate é executado antes de inserir um novo usuário.
// Gera automaticamente o UUID se não foi definido.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// Volunteer representa dados adicionais de um voluntário.
type Volunteer struct {
	UserID         uuid.UUID `gorm:"type:uniqueidentifier;primaryKey" json:"user_id"`
	DedicatedHours float64   `gorm:"default:0" json:"dedicated_hours"`
	IsVerified     bool      `gorm:"default:false" json:"is_verified"`
	RatingAvg      float64   `gorm:"default:0" json:"rating_avg"`
	RatingCount    int       `gorm:"default:0" json:"rating_count"`

	// Relacionamento com User
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName define o nome da tabela no banco de dados.
func (Volunteer) TableName() string {
	return "volunteers"
}

// Elderly representa dados adicionais de um idoso.
type Elderly struct {
	UserID           uuid.UUID `gorm:"type:uniqueidentifier;primaryKey" json:"user_id"`
	EmergencyContact string    `gorm:"size:255" json:"emergency_contact,omitempty"`
	NeedsAssistance  bool      `gorm:"default:false" json:"needs_assistance"`

	// Relacionamento com User
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName define o nome da tabela no banco de dados.
func (Elderly) TableName() string {
	return "elderly"
}

// Institution representa dados adicionais de uma instituição.
type Institution struct {
	UserID          uuid.UUID `gorm:"type:uniqueidentifier;primaryKey" json:"user_id"`
	InstitutionType string    `gorm:"size:100" json:"institution_type"`
	VisitDays       string    `gorm:"size:255" json:"visit_days,omitempty"`
	VisitTime       string    `gorm:"size:50" json:"visit_time,omitempty"`
	VisitDuration   string    `gorm:"size:50" json:"visit_duration,omitempty"`
	ResponsibleName string    `gorm:"size:255" json:"responsible_name,omitempty"`

	// Relacionamento com User
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName define o nome da tabela no banco de dados.
func (Institution) TableName() string {
	return "institutions"
}
