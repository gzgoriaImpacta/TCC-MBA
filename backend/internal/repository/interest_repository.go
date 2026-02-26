// Package repository contém as implementações de acesso a dados.
package repository

import (
	"errors"

	"amigos-terceira-idade/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InterestRepository gerencia as operações de banco de dados para interesses.
type InterestRepository struct {
	db *gorm.DB
}

// NewInterestRepository cria uma nova instância do repositório de interesses.
func NewInterestRepository(db *gorm.DB) *InterestRepository {
	return &InterestRepository{db: db}
}

// Create insere um novo interesse no banco de dados.
func (r *InterestRepository) Create(interest *domain.Interest) error {
	return r.db.Create(interest).Error
}

// FindAll retorna todos os interesses disponíveis.
func (r *InterestRepository) FindAll() ([]domain.Interest, error) {
	var interests []domain.Interest
	err := r.db.Order("name ASC").Find(&interests).Error
	if err != nil {
		return nil, err
	}
	return interests, nil
}

// FindByID busca um interesse pelo ID.
func (r *InterestRepository) FindByID(id uuid.UUID) (*domain.Interest, error) {
	var interest domain.Interest
	err := r.db.First(&interest, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("interesse não encontrado")
		}
		return nil, err
	}
	return &interest, nil
}

// FindByIDs busca múltiplos interesses pelos IDs.
// Útil para vincular interesses a um usuário no cadastro.
func (r *InterestRepository) FindByIDs(ids []uuid.UUID) ([]domain.Interest, error) {
	var interests []domain.Interest
	err := r.db.Where("id IN ?", ids).Find(&interests).Error
	if err != nil {
		return nil, err
	}
	return interests, nil
}

// FindByName busca um interesse pelo nome.
func (r *InterestRepository) FindByName(name string) (*domain.Interest, error) {
	var interest domain.Interest
	err := r.db.First(&interest, "name = ?", name).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("interesse não encontrado")
		}
		return nil, err
	}
	return &interest, nil
}

// SeedDefaults insere os interesses padrão do sistema.
// Deve ser executado na inicialização do banco.
func (r *InterestRepository) SeedDefaults() error {
	defaults := domain.DefaultInterests()
	for _, interest := range defaults {
		// Verifica se já existe antes de inserir
		var existing domain.Interest
		err := r.db.First(&existing, "name = ?", interest.Name).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := r.db.Create(&interest).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
