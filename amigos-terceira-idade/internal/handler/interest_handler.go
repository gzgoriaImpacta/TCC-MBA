// Package handler contém os handlers HTTP da aplicação.
package handler

import (
	"net/http"

	"amigos-terceira-idade/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// InterestHandler gerencia os endpoints de interesses.
type InterestHandler struct {
	interestService *service.InterestService
}

// NewInterestHandler cria uma nova instância do handler de interesses.
func NewInterestHandler(interestService *service.InterestService) *InterestHandler {
	return &InterestHandler{
		interestService: interestService,
	}
}

// GetAll godoc
// @Summary Lista todos os interesses
// @Description Retorna todos os interesses disponíveis no sistema
// @Tags Interests
// @Produce json
// @Success 200 {object} Response
// @Router /interests [get]
func (h *InterestHandler) GetAll(c *gin.Context) {
	interests, err := h.interestService.GetAll()
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "FETCH_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, interests)
}

// GetByID godoc
// @Summary Busca um interesse pelo ID
// @Description Retorna os dados de um interesse específico
// @Tags Interests
// @Produce json
// @Param id path string true "ID do interesse"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Router /interests/{id} [get]
func (h *InterestHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID inválido")
		return
	}

	interest, err := h.interestService.GetByID(id)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "INTEREST_NOT_FOUND", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, interest)
}
