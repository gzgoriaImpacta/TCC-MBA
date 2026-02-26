// Package handler contém os handlers HTTP da aplicação.
package handler

import (
	"net/http"

	"amigos-terceira-idade/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler gerencia os endpoints de autenticação.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler cria uma nova instância do handler de autenticação.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Cadastra um novo usuário
// @Description Cria uma nova conta de voluntário, idoso ou instituição
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "Dados do cadastro"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Dados inválidos: "+err.Error())
		return
	}

	// Valida o tipo de usuário
	validTypes := map[string]bool{"VOLUNTEER": true, "ELDERLY": true, "INSTITUTION": true}
	if !validTypes[req.UserType] {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_USER_TYPE", "Tipo de usuário inválido. Use: VOLUNTEER, ELDERLY ou INSTITUTION")
		return
	}

	result, err := h.authService.Register(req)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "REGISTER_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusCreated, result)
}

// Login godoc
// @Summary Realiza login
// @Description Autentica um usuário e retorna os tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "Credenciais"
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Dados inválidos: "+err.Error())
		return
	}

	result, err := h.authService.Login(req)
	if err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "LOGIN_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, result)
}

// RefreshTokenRequest contém o refresh token para renovação.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken godoc
// @Summary Renova os tokens
// @Description Gera novos tokens a partir do refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Dados inválidos: "+err.Error())
		return
	}

	result, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "REFRESH_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, result)
}
