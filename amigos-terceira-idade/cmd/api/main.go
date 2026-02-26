// Package main é o ponto de entrada da aplicação.
// Inicializa as dependências e inicia o servidor HTTP.
package main

import (
	"log"

	"amigos-terceira-idade/internal/config"
	"amigos-terceira-idade/internal/domain"
	"amigos-terceira-idade/internal/handler"
	"amigos-terceira-idade/internal/repository"
	"amigos-terceira-idade/internal/service"
	"amigos-terceira-idade/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Carrega as configurações
	cfg := config.Load()

	// Configura o modo do Gin
	gin.SetMode(cfg.Server.Mode)

	// Conecta ao banco de dados
	db, err := database.NewConnection(database.DatabaseConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Database: cfg.Database.Database,
	})
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer database.Close(db)

	// Executa as migrations automaticamente
	log.Println("Executando migrations...")
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Interest{},
		&domain.Volunteer{},
		&domain.Elderly{},
		&domain.Institution{},
		&domain.Connection{},
		&domain.Appointment{},
	)
	if err != nil {
		log.Fatalf("Erro ao executar migrations: %v", err)
	}
	log.Println("Migrations executadas com sucesso")

	// Inicializa os repositórios
	userRepo := repository.NewUserRepository(db)
	interestRepo := repository.NewInterestRepository(db)
	connectionRepo := repository.NewConnectionRepository(db)
	appointmentRepo := repository.NewAppointmentRepository(db)

	// Insere os interesses padrão
	log.Println("Inserindo interesses padrão...")
	if err := interestRepo.SeedDefaults(); err != nil {
		log.Printf("Aviso: erro ao inserir interesses padrão: %v", err)
	}

	// Inicializa os serviços
	authService := service.NewAuthService(userRepo, interestRepo, cfg.JWT)
	userService := service.NewUserService(userRepo, interestRepo)
	interestService := service.NewInterestService(interestRepo)
	matchingService := service.NewMatchingService(userRepo, connectionRepo)
	appointmentService := service.NewAppointmentService(appointmentRepo, userRepo)

	// Inicializa os handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	interestHandler := handler.NewInterestHandler(interestService)
	matchingHandler := handler.NewMatchingHandler(matchingService)
	appointmentHandler := handler.NewAppointmentHandler(appointmentService)

	// Configura o router
	router := handler.NewRouter(
		authHandler,
		userHandler,
		interestHandler,
		matchingHandler,
		appointmentHandler,
		authService,
	)

	// Cria o engine do Gin
	engine := gin.Default()

	// Configura as rotas
	router.Setup(engine)

	// Inicia o servidor
	log.Printf("Servidor iniciando na porta %s...", cfg.Server.Port)
	if err := engine.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
