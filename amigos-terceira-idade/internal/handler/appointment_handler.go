// Package handler contém os handlers HTTP da aplicação.
package handler

import (
	"net/http"

	"amigos-terceira-idade/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AppointmentHandler gerencia os endpoints de agendamentos.
type AppointmentHandler struct {
	appointmentService *service.AppointmentService
}

// NewAppointmentHandler cria uma nova instância do handler de agendamentos.
func NewAppointmentHandler(appointmentService *service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: appointmentService,
	}
}

// Create godoc
// @Summary Cria um agendamento
// @Description Voluntário envia convite para conversa
// @Tags Appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateAppointmentRequest true "Dados do agendamento"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Router /appointments [post]
func (h *AppointmentHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req service.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Dados inválidos: "+err.Error())
		return
	}

	appointment, err := h.appointmentService.Create(userID, req)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "CREATE_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusCreated, appointment)
}

// GetMy godoc
// @Summary Lista meus agendamentos
// @Description Retorna todos os agendamentos do usuário
// @Tags Appointments
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response
// @Router /appointments [get]
func (h *AppointmentHandler) GetMy(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	appointments, err := h.appointmentService.GetMyAppointments(userID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "FETCH_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, appointments)
}

// GetUpcoming godoc
// @Summary Lista próximos agendamentos
// @Description Retorna agendamentos futuros confirmados
// @Tags Appointments
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response
// @Router /appointments/upcoming [get]
func (h *AppointmentHandler) GetUpcoming(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	appointments, err := h.appointmentService.GetUpcoming(userID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "FETCH_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, appointments)
}

// GetByID godoc
// @Summary Busca um agendamento pelo ID
// @Description Retorna os detalhes de um agendamento
// @Tags Appointments
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do agendamento"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Router /appointments/{id} [get]
func (h *AppointmentHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID inválido")
		return
	}

	appointment, err := h.appointmentService.GetByID(id)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "APPOINTMENT_NOT_FOUND", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, appointment)
}

// GetReceivedInvitations godoc
// @Summary Lista convites recebidos
// @Description Retorna convites pendentes recebidos
// @Tags Invitations
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response
// @Router /invitations/received [get]
func (h *AppointmentHandler) GetReceivedInvitations(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	invitations, err := h.appointmentService.GetReceivedInvitations(userID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "FETCH_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, invitations)
}

// GetSentInvitations godoc
// @Summary Lista convites enviados
// @Description Retorna convites pendentes enviados
// @Tags Invitations
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response
// @Router /invitations/sent [get]
func (h *AppointmentHandler) GetSentInvitations(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	invitations, err := h.appointmentService.GetSentInvitations(userID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "FETCH_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, invitations)
}

// Accept godoc
// @Summary Aceita um convite
// @Description Idoso/instituição aceita convite de conversa
// @Tags Appointments
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do agendamento"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /appointments/{id}/accept [post]
func (h *AppointmentHandler) Accept(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID inválido")
		return
	}

	if err := h.appointmentService.Accept(id, userID); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "ACCEPT_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"message": "Convite aceito com sucesso"})
}

// Decline godoc
// @Summary Recusa um convite
// @Description Idoso/instituição recusa convite de conversa
// @Tags Appointments
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do agendamento"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /appointments/{id}/decline [post]
func (h *AppointmentHandler) Decline(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID inválido")
		return
	}

	if err := h.appointmentService.Decline(id, userID); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "DECLINE_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"message": "Convite recusado"})
}

// Cancel godoc
// @Summary Cancela um agendamento
// @Description Cancela um agendamento existente
// @Tags Appointments
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do agendamento"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /appointments/{id} [delete]
func (h *AppointmentHandler) Cancel(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID inválido")
		return
	}

	if err := h.appointmentService.Cancel(id, userID); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "CANCEL_ERROR", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"message": "Agendamento cancelado"})
}
