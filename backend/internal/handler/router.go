// Package handler contém os handlers HTTP da aplicação.
package handler

import (
	"amigos-terceira-idade/internal/middleware"
	"amigos-terceira-idade/internal/service"
	"github.com/gin-gonic/gin"
)

// Router configura todas as rotas da API.
type Router struct {
	authHandler        *AuthHandler
	userHandler        *UserHandler
	interestHandler    *InterestHandler
	matchingHandler    *MatchingHandler
	appointmentHandler *AppointmentHandler
	authService        *service.AuthService
}

// NewRouter cria uma nova instância do router.
func NewRouter(
	authHandler *AuthHandler,
	userHandler *UserHandler,
	interestHandler *InterestHandler,
	matchingHandler *MatchingHandler,
	appointmentHandler *AppointmentHandler,
	authService *service.AuthService,
) *Router {
	return &Router{
		authHandler:        authHandler,
		userHandler:        userHandler,
		interestHandler:    interestHandler,
		matchingHandler:    matchingHandler,
		appointmentHandler: appointmentHandler,
		authService:        authService,
	}
}

// Setup configura todas as rotas no engine do Gin.
func (r *Router) Setup(engine *gin.Engine) {
	// Middleware global de CORS
	engine.Use(middleware.CORSMiddleware())

	// Grupo base da API
	api := engine.Group("/api/v1")

	// Rotas públicas (sem autenticação)
	r.setupPublicRoutes(api)

	// Rotas protegidas (com autenticação)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(r.authService))
	r.setupProtectedRoutes(protected)
}

// setupPublicRoutes configura as rotas que não precisam de autenticação.
func (r *Router) setupPublicRoutes(api *gin.RouterGroup) {
	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "amigos-terceira-idade"})
	})

	// Autenticação
	auth := api.Group("/auth")
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/refresh", r.authHandler.RefreshToken)
	}

	// Interesses (público para mostrar no cadastro)
	api.GET("/interests", r.interestHandler.GetAll)
	api.GET("/interests/:id", r.interestHandler.GetByID)
}

// setupProtectedRoutes configura as rotas que precisam de autenticação.
func (r *Router) setupProtectedRoutes(api *gin.RouterGroup) {
	// Usuários
	users := api.Group("/users")
	{
		users.GET("/me", r.userHandler.GetMe)
		users.PUT("/me", r.userHandler.UpdateMe)
		users.DELETE("/me", r.userHandler.Deactivate)
		users.GET("/:id", r.userHandler.GetByID)
	}

	// Pareamento
	matching := api.Group("/matching")
	{
		matching.GET("/suggestions", r.matchingHandler.GetSuggestions)
		matching.POST("/connect", r.matchingHandler.Connect)
		matching.GET("/connections", r.matchingHandler.GetConnections)
		matching.POST("/connections/:id/accept", r.matchingHandler.AcceptConnection)
		matching.POST("/connections/:id/reject", r.matchingHandler.RejectConnection)
	}

	// Agendamentos
	appointments := api.Group("/appointments")
	{
		appointments.POST("", r.appointmentHandler.Create)
		appointments.GET("", r.appointmentHandler.GetMy)
		appointments.GET("/upcoming", r.appointmentHandler.GetUpcoming)
		appointments.GET("/:id", r.appointmentHandler.GetByID)
		appointments.POST("/:id/accept", r.appointmentHandler.Accept)
		appointments.POST("/:id/decline", r.appointmentHandler.Decline)
		appointments.DELETE("/:id", r.appointmentHandler.Cancel)
	}

	// Convites (atalhos para agendamentos pendentes)
	invitations := api.Group("/invitations")
	{
		invitations.GET("/received", r.appointmentHandler.GetReceivedInvitations)
		invitations.GET("/sent", r.appointmentHandler.GetSentInvitations)
	}
}
