// Package domain contém as entidades de negócio da aplicação.
package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Interest representa um interesse/hobby que usuários podem ter.
// Usado para fazer o pareamento entre voluntários e idosos.
type Interest struct {
	ID        uuid.UUID `gorm:"type:uniqueidentifier;default:NEWID();primaryKey"`
	Name      string    `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Icon      string    `gorm:"size:50" json:"icon,omitempty"` // Emoji ou nome do ícone
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName define o nome da tabela no banco de dados.
func (Interest) TableName() string {
	return "interests"
}

// BeforeCreate é executado antes de inserir um novo interesse.
func (i *Interest) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// DefaultInterests retorna a lista de interesses padrão do sistema.
// Estes são criados automaticamente na inicialização do banco.
func DefaultInterests() []Interest {
	return []Interest{
		{Name: "Instrumentos musicais", Icon: "🎸"},
		{Name: "Jogos de tabuleiro", Icon: "🎲"},
		{Name: "Caminhadas", Icon: "🚶"},
		{Name: "Leitura", Icon: "📚"},
		{Name: "Palavras cruzadas", Icon: "✏️"},
		{Name: "Música", Icon: "🎵"},
		{Name: "Xadrez", Icon: "♟️"},
		{Name: "Jardinagem", Icon: "🌱"},
		{Name: "Atividades manuais", Icon: "🎨"},
		{Name: "Conversa em grupo", Icon: "💬"},
	}
}
