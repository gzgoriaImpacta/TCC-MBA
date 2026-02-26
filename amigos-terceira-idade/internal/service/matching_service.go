// Package service contém a lógica de negócio da aplicação.
package service

import (
	"errors"

	"amigos-terceira-idade/internal/domain"
	"amigos-terceira-idade/internal/repository"

	"github.com/google/uuid"
)

// MatchingService gerencia o pareamento entre voluntários e idosos/instituições.
type MatchingService struct {
	userRepo       repository.UserRepositoryInterface
	connectionRepo repository.ConnectionRepositoryInterface
}

// NewMatchingService cria uma nova instância do serviço de pareamento.
func NewMatchingService(
	userRepo repository.UserRepositoryInterface,
	connectionRepo repository.ConnectionRepositoryInterface,
) *MatchingService {
	return &MatchingService{
		userRepo:       userRepo,
		connectionRepo: connectionRepo,
	}
}

// MatchSuggestion representa uma sugestão de pareamento com score de compatibilidade.
type MatchSuggestion struct {
	User             domain.User `json:"user"`
	MatchedInterests int         `json:"matched_interests"` // Quantidade de interesses em comum
	MatchScore       float64     `json:"match_score"`       // Porcentagem de match (0-100)
}

// GetSuggestions retorna sugestões de pareamento para um voluntário.
// Ordena por quantidade de interesses em comum.
func (s *MatchingService) GetSuggestions(volunteerID uuid.UUID, filterType string) ([]MatchSuggestion, error) {
	// Busca o voluntário para obter seus interesses
	volunteer, err := s.userRepo.FindByID(volunteerID)
	if err != nil {
		return nil, err
	}

	if volunteer.UserType != domain.UserTypeVolunteer {
		return nil, errors.New("apenas voluntários podem buscar conexões")
	}

	// Define o tipo de usuário a buscar
	var targetTypes []domain.UserType
	switch filterType {
	case "elderly":
		targetTypes = []domain.UserType{domain.UserTypeElderly}
	case "institution":
		targetTypes = []domain.UserType{domain.UserTypeInstitution}
	default:
		targetTypes = []domain.UserType{domain.UserTypeElderly, domain.UserTypeInstitution}
	}

	// Busca os usuários dos tipos especificados
	var suggestions []MatchSuggestion
	for _, userType := range targetTypes {
		users, err := s.userRepo.FindByType(userType)
		if err != nil {
			return nil, err
		}

		for _, user := range users {
			// Verifica se já existe conexão
			exists, _ := s.connectionRepo.Exists(volunteerID, user.ID)
			if exists {
				continue // Pula usuários já conectados
			}

			// Calcula a quantidade de interesses em comum
			matchedCount := s.countMatchedInterests(volunteer.Interests, user.Interests)

			// Calcula o score de compatibilidade
			var matchScore float64
			if len(volunteer.Interests) > 0 {
				matchScore = float64(matchedCount) / float64(len(volunteer.Interests)) * 100
			}

			suggestions = append(suggestions, MatchSuggestion{
				User:             user,
				MatchedInterests: matchedCount,
				MatchScore:       matchScore,
			})
		}
	}

	// Ordena por quantidade de matches (maior primeiro)
	s.sortByMatchScore(suggestions)

	return suggestions, nil
}

// countMatchedInterests conta quantos interesses são comuns entre duas listas.
func (s *MatchingService) countMatchedInterests(interests1, interests2 []domain.Interest) int {
	interestMap := make(map[uuid.UUID]bool)
	for _, i := range interests1 {
		interestMap[i.ID] = true
	}

	count := 0
	for _, i := range interests2 {
		if interestMap[i.ID] {
			count++
		}
	}
	return count
}

// sortByMatchScore ordena as sugestões por score de compatibilidade (maior primeiro).
func (s *MatchingService) sortByMatchScore(suggestions []MatchSuggestion) {
	// Bubble sort simples - adequado para listas pequenas
	for i := 0; i < len(suggestions)-1; i++ {
		for j := 0; j < len(suggestions)-i-1; j++ {
			if suggestions[j].MatchScore < suggestions[j+1].MatchScore {
				suggestions[j], suggestions[j+1] = suggestions[j+1], suggestions[j]
			}
		}
	}
}

// Connect cria uma nova conexão entre um voluntário e um idoso/instituição.
func (s *MatchingService) Connect(volunteerID, targetID uuid.UUID) (*domain.Connection, error) {
	// Verifica se já existe conexão
	exists, err := s.connectionRepo.Exists(volunteerID, targetID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("conexão já existe")
	}

	// Busca os usuários para validação
	volunteer, err := s.userRepo.FindByID(volunteerID)
	if err != nil {
		return nil, err
	}
	if volunteer.UserType != domain.UserTypeVolunteer {
		return nil, errors.New("apenas voluntários podem iniciar conexões")
	}

	target, err := s.userRepo.FindByID(targetID)
	if err != nil {
		return nil, err
	}
	if target.UserType == domain.UserTypeVolunteer {
		return nil, errors.New("não é possível conectar com outro voluntário")
	}

	// Calcula interesses em comum
	matchedCount := s.countMatchedInterests(volunteer.Interests, target.Interests)

	// Cria a conexão
	connection := &domain.Connection{
		ID:               uuid.New(),
		VolunteerID:      volunteerID,
		TargetID:         targetID,
		TargetType:       target.UserType,
		Status:           domain.ConnectionStatusPending,
		MatchedInterests: matchedCount,
	}

	if err := s.connectionRepo.Create(connection); err != nil {
		return nil, err
	}

	return connection, nil
}

// GetConnections retorna as conexões de um usuário.
func (s *MatchingService) GetConnections(userID uuid.UUID) ([]domain.Connection, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user.UserType == domain.UserTypeVolunteer {
		return s.connectionRepo.FindByVolunteerID(userID)
	}
	return s.connectionRepo.FindByTargetID(userID)
}

// AcceptConnection aceita uma conexão pendente.
func (s *MatchingService) AcceptConnection(connectionID uuid.UUID) error {
	return s.connectionRepo.UpdateStatus(connectionID, domain.ConnectionStatusAccepted)
}

// RejectConnection rejeita uma conexão pendente.
func (s *MatchingService) RejectConnection(connectionID uuid.UUID) error {
	return s.connectionRepo.UpdateStatus(connectionID, domain.ConnectionStatusRejected)
}
