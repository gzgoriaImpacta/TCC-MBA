// Package service contém a lógica de negócio da aplicação.
package service

import (
	"amigos-terceira-idade/internal/domain"
	"amigos-terceira-idade/internal/repository"
	"github.com/google/uuid"
)

// InterestService gerencia os interesses disponíveis no sistema.
type InterestService struct {
	interestRepo repository.InterestRepositoryInterface
}

// NewInterestService cria uma nova instância do serviço de interesses.
func NewInterestService(interestRepo repository.InterestRepositoryInterface) *InterestService {
	return &InterestService{
		interestRepo: interestRepo,
	}
}

// GetAll retorna todos os interesses disponíveis.
func (s *InterestService) GetAll() ([]domain.Interest, error) {
	return s.interestRepo.FindAll()
}

// GetByID busca um interesse pelo ID.
func (s *InterestService) GetByID(id uuid.UUID) (*domain.Interest, error) {
	return s.interestRepo.FindByID(id)
}

// SeedDefaults insere os interesses padrão do sistema.
func (s *InterestService) SeedDefaults() error {
	return s.interestRepo.SeedDefaults()
}
