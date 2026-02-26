// Package database contém a configuração de conexão com o banco de dados.
// Utiliza GORM como ORM para facilitar as operações com SQL Server.
package database

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig contém os parâmetros de conexão com o banco.
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// NewConnection cria uma nova conexão com o SQL Server.
// Retorna uma instância do GORM configurada e pronta para uso.
func NewConnection(cfg DatabaseConfig) (*gorm.DB, error) {
	// Monta a string de conexão no formato SQL Server
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	// Configura o GORM com logging habilitado para debug
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com SQL Server: %w", err)
	}

	log.Println("Conexão com SQL Server estabelecida com sucesso")
	return db, nil
}

// Close fecha a conexão com o banco de dados.
// Deve ser chamado ao encerrar a aplicação.
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("erro ao obter conexão SQL: %w", err)
	}
	return sqlDB.Close()
}
