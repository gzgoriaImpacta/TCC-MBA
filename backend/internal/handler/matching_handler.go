// Package handler contém os handlers HTTP da aplicação.
package handler

import (
	"net/http"

	"amigos-terceira-idade/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MatchingHandler gerencia os endpoints de pareamento.
type MatchingHandler struct {
	matchingService *service.MatchingService
}

// NewMatchingHandler cria uma nova instância do handler de pareamento.
func NewMatchingHandler(matchingService *service.MatchingService) *MatchingHandler {
	return &MatchingHandler{
		matchingService: matchingService,
	}
}

// GetSuggestions godoc
// @Summary Retorna sugestões de pareamento
// @Description Lista idosos/instituições compatíveis com o voluntário
// @Tags Matching
// @Produce json
// @Security BearerAuth
// @Param type query string false "Filtrar por tipo: elderly ou institution"
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Router /matching/suggestions [get]
func (h *MatchingHandler) GetSuggestions(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	filterType := c.Query("type")

	suggestions, err := h.matchingService.GetSuggestions(userID, filterType)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "SUGGESTIONS_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, suggestions)
}

// ConnectRequest contém os dados para criar uma conexão.
type ConnectRequest struct {
	TargetID uuid.UUID `json:"target_id" binding:"required"`
}

// Connect godoc
// @Summary Cria uma conexão
// @Description Voluntário solicita conexão com idoso/instituição
// @Tags Matching
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ConnectRequest true "ID do alvo"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Router /matching/connect [post]
func (h *MatchingHandler) Connect(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req ConnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Dados inválidos: "+err.Error())
		return
	}

	connection, err := h.matchingService.Connect(userID, req.TargetID)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "CONNECT_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusCreated, connection)
}

// GetConnections godoc
// @Summary Lista as conexões do usuário
// @Description Retorna todas as conexões do usuário autenticado
// @Tags Matching
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response
// @Router /matching/connections [get]
func (h *MatchingHandler) GetConnections(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	connections, err := h.matchingService.GetConnections(userID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "FETCH_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, connections)
}

// AcceptConnection godoc
// @Summary Aceita uma conexão
// @Description Idoso/instituição aceita conexão de um voluntário
// @Tags Matching
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da conexão"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /matching/connections/{id}/accept [post]
func (h *MatchingHandler) AcceptConnection(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID inválido")
		return
	}

	if err := h.matchingService.AcceptConnection(id); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "ACCEPT_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"message": "Conexão aceita com sucesso"})
}

// RejectConnection godoc
// @Summary Rejeita uma conexão
// @Description Idoso/instituição rejeita conexão de um voluntário
// @Tags Matching
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da conexão"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /matching/connections/{id}/reject [post]
func (h *MatchingHandler) RejectConnection(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID inválido")
		return
	}

	if err := h.matchingService.RejectConnection(id); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "REJECT_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"message": "Conexão rejeitada"})
}
