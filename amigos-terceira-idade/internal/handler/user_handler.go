// Package handler contém os handlers HTTP da aplicação.
package handler

import (
	"net/http"

	"amigos-terceira-idade/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler gerencia os endpoints de usuários.
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler cria uma nova instância do handler de usuários.
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetMe godoc
// @Summary Retorna o perfil do usuário autenticado
// @Description Retorna os dados do usuário logado
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Router /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	user, err := h.userService.GetByID(userID)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "USER_NOT_FOUND", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, user)
}

// UpdateMe godoc
// @Summary Atualiza o perfil do usuário autenticado
// @Description Atualiza os dados do usuário logado
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.UpdateProfileRequest true "Dados para atualização"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /users/me [put]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Dados inválidos: "+err.Error())
		return
	}

	user, err := h.userService.UpdateProfile(userID, req)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "UPDATE_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, user)
}

// GetByID godoc
// @Summary Retorna um usuário pelo ID
// @Description Retorna os dados públicos de um usuário
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do usuário"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID inválido")
		return
	}

	user, err := h.userService.GetByID(id)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "USER_NOT_FOUND", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, user)
}

// Deactivate godoc
// @Summary Desativa a conta do usuário
// @Description Desativa a conta do usuário autenticado
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /users/me [delete]
func (h *UserHandler) Deactivate(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	if err := h.userService.Deactivate(userID); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "DEACTIVATE_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"message": "Conta desativada com sucesso"})
}
