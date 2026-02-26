// Package service contém a lógica de negócio da aplicação.
package service

import (
	"amigos-terceira-idade/internal/domain"
	"amigos-terceira-idade/internal/repository"
	"fmt"

	"github.com/google/uuid"
)

// UserService gerencia as operações relacionadas a usuários.
type UserService struct {
	userRepo     repository.UserRepositoryInterface
	interestRepo repository.InterestRepositoryInterface
}

// NewUserService cria uma nova instância do serviço de usuários.
func NewUserService(
	userRepo repository.UserRepositoryInterface,
	interestRepo repository.InterestRepositoryInterface,
) *UserService {
	return &UserService{
		userRepo:     userRepo,
		interestRepo: interestRepo,
	}
}

// UpdateProfileRequest contém os dados para atualização de perfil.
type UpdateProfileRequest struct {
	Name        string      `json:"name"`
	Age         int         `json:"age"`
	Bio         string      `json:"bio"`
	Phone       string      `json:"phone"`
	PhotoURL    string      `json:"photo_url"`
	InterestIDs []uuid.UUID `json:"interest_ids"`
}

// GetByID busca um usuário pelo ID.
func (s *UserService) GetByID(id uuid.UUID) (*domain.User, error) {
	return s.userRepo.FindByID(id)
}

// UpdateProfile atualiza o perfil de um usuário.
func (s *UserService) UpdateProfile(userID uuid.UUID, req UpdateProfileRequest) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// Atualiza apenas os campos fornecidos
	if req.Name != "" {
		user.Name = req.Name
	}
	// if req.Age > 0 {
	// 	user.Age = req.Age
	// }
	fmt.Println(req)
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	// // if req.PhotoURL != "" {
	// 	user.PhotoURL = req.PhotoURL
	// }

	// Atualiza os interesses se foram fornecidos
	if len(req.InterestIDs) > 0 {
		interests, err := s.interestRepo.FindByIDs(req.InterestIDs)
		if err != nil {
			return nil, err
		}
		if err := s.userRepo.UpdateInterests(userID, interests); err != nil {
			return nil, err
		}
	}

	// Salva as alterações
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	updatedUser, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	fmt.Println("teste")
	fmt.Println(updatedUser)

	// Retorna o usuário atualizado com interesses
	return updatedUser, nil
}

// ListByType retorna todos os usuários de um determinado tipo.
func (s *UserService) ListByType(userType domain.UserType) ([]domain.User, error) {
	return s.userRepo.FindByType(userType)
}

// Deactivate desativa um usuário (soft delete).
func (s *UserService) Deactivate(userID uuid.UUID) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	user.IsActive = false
	return s.userRepo.Update(user)
}
