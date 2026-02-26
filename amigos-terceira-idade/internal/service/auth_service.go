// Package service contém a lógica de negócio da aplicação.
// Os services orquestram as operações entre repositórios e aplicam regras de negócio.
package service

import (
	"errors"
	"fmt"
	"time"

	"amigos-terceira-idade/internal/config"
	"amigos-terceira-idade/internal/domain"
	"amigos-terceira-idade/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

// AuthService gerencia a autenticação e autorização de usuários.
type AuthService struct {
	userRepo     repository.UserRepositoryInterface
	interestRepo repository.InterestRepositoryInterface
	jwtConfig    config.JWTConfig
}

// NewAuthService cria uma nova instância do serviço de autenticação.
func NewAuthService(
	userRepo repository.UserRepositoryInterface,
	interestRepo repository.InterestRepositoryInterface,
	jwtConfig config.JWTConfig,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		interestRepo: interestRepo,
		jwtConfig:    jwtConfig,
	}
}

// RegisterRequest contém os dados necessários para cadastro.
type RegisterRequest struct {
	Name        string      `json:"name" binding:"required"`
	Email       string      `json:"email" binding:"required,email"`
	Password    string      `json:"password" binding:"required,min=6"`
	Age         int         `json:"age"`
	Bio         string      `json:"bio"`
	UserType    string      `json:"user_type" binding:"required"` // VOLUNTEER, ELDERLY, INSTITUTION
	InterestIDs []uuid.UUID `json:"interest_ids"`
}

// LoginRequest contém os dados necessários para login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse contém os tokens de autenticação.
type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *domain.User `json:"user"`
}

// TokenClaims representa os dados contidos no JWT.
type TokenClaims struct {
	UserID   uuid.UUID       `json:"user_id"`
	Email    string          `json:"email"`
	UserType domain.UserType `json:"user_type"`
	jwt.RegisteredClaims
}

// Register realiza o cadastro de um novo usuário.
func (s *AuthService) Register(req RegisterRequest) (*AuthResponse, error) {
	// Verifica se o email já está em uso
	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email já está em uso")
	}

	// Gera o hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("erro ao processar senha")
	}

	// Cria o usuário
	user := &domain.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Age:          req.Age,
		Bio:          req.Bio,
		UserType:     domain.UserType(req.UserType),
		IsActive:     true,
	}

	// Busca os interesses selecionados
	if len(req.InterestIDs) > 0 {
		fmt.Println("Interesses ID: ", req.InterestIDs)
		interests, err := s.interestRepo.FindByIDs(req.InterestIDs)
		for i := range interests {
			interests[i].ID = FixUUID(interests[i].ID)
		}

		if err != nil {
			return nil, err
		}
		user.Interests = interests
	}

	// Salva o usuário no banco
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Gera os tokens
	return s.generateAuthResponse(user)
}

// Login realiza a autenticação de um usuário.
func (s *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	// Busca o usuário pelo email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	// Verifica a senha
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	// Verifica se o usuário está ativo
	if !user.IsActive {
		return nil, errors.New("usuário desativado")
	}

	// Gera os tokens
	return s.generateAuthResponse(user)
}

// ValidateToken valida um token JWT e retorna as claims.
func (s *AuthService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtConfig.SecretKey), nil
	})

	if err != nil {
		return nil, errors.New("token inválido")
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}

// RefreshToken gera novos tokens a partir de um refresh token válido.
func (s *AuthService) RefreshToken(refreshToken string) (*AuthResponse, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	return s.generateAuthResponse(user)
}

// generateAuthResponse gera os tokens de autenticação para um usuário.
func (s *AuthService) generateAuthResponse(user *domain.User) (*AuthResponse, error) {
	accessToken, err := s.generateToken(user, time.Duration(s.jwtConfig.AccessTokenExpiry)*time.Hour)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user, time.Duration(s.jwtConfig.RefreshTokenExpiry)*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// generateToken gera um token JWT com o tempo de expiração especificado.
func (s *AuthService) generateToken(user *domain.User, expiry time.Duration) (string, error) {
	claims := TokenClaims{
		UserID:   user.ID,
		Email:    user.Email,
		UserType: user.UserType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtConfig.SecretKey))
}
