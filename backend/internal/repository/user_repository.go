// Package repository contém as implementações de acesso a dados.
// Cada repositório é responsável por uma entidade específica.
package repository

import (
	"amigos-terceira-idade/internal/domain"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func FixUUID(u uuid.UUID) uuid.UUID {
	b := u
	// inverter os 4 primeiros bytes
	b[0], b[1], b[2], b[3] = u[3], u[2], u[1], u[0]
	// inverter os próximos 2 bytes
	b[4], b[5] = u[5], u[4]
	// inverter os próximos 2 bytes
	b[6], b[7] = u[7], u[6]
	// os últimos 8 bytes permanecem iguais
	return b
}

// UserRepository gerencia as operações de banco de dados para usuários.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository cria uma nova instância do repositório de usuários.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create insere um novo usuário no banco de dados.
func (r *UserRepository) Create(user *domain.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Omit("Interests.*").Create(user).Error; err != nil {
			return err
		}
		if len(user.Interests) > 0 {
			if err := tx.Model(user).
				Association("Interests").
				Replace(user.Interests); err != nil {
				return err
			}
		}

		return nil
	})
	// return r.db.Create(user).Error
}

// FindByID busca um usuário pelo ID.
// Retorna erro se não encontrar.
func (r *UserRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User

	err := r.db.
		Debug().
		Preload("Interests").
		First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}

	return &user, nil
}

// FindByEmail busca um usuário pelo email.
// Usado principalmente no login.
func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Interests").First(&user, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}
	// Transformar UUID em string canônica
	user.ID = FixUUID(user.ID)

	return &user, nil
}

// Update atualiza os dados de um usuário existente.
func (r *UserRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// Delete remove um usuário do banco de dados (soft delete recomendado).
func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.User{}, "id = ?", id).Error
}

// ExistsByEmail verifica se já existe um usuário com o email informado.
// Usado na validação de cadastro.
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindByType busca todos os usuários de um determinado tipo.
// Útil para listar todos os voluntários, idosos ou instituições.
func (r *UserRepository) FindByType(userType domain.UserType) ([]domain.User, error) {
	var users []domain.User
	err := r.db.Preload("Interests").
		Where("user_type = ? AND is_active = ?", userType, true).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// AddInterests adiciona interesses a um usuário.
func (r *UserRepository) AddInterests(userID uuid.UUID, interests []domain.Interest) error {
	user, err := r.FindByID(userID)
	if err != nil {
		return err
	}
	return r.db.Model(user).Association("Interests").Append(interests)
}

// RemoveInterest remove um interesse de um usuário.
func (r *UserRepository) RemoveInterest(userID uuid.UUID, interestID uuid.UUID) error {
	user, err := r.FindByID(userID)
	if err != nil {
		return err
	}
	interest := domain.Interest{ID: interestID}
	return r.db.Model(user).Association("Interests").Delete(&interest)
}

// UpdateInterests substitui todos os interesses de um usuário.
func (r *UserRepository) UpdateInterests(userID uuid.UUID, interests []domain.Interest) error {
	user, err := r.FindByID(userID)
	if err != nil {
		return err
	}
	return r.db.Model(user).Association("Interests").Replace(interests)
}
