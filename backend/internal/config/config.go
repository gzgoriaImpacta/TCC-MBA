// Package config contém as configurações da aplicação.
// Todas as variáveis de ambiente e configurações são centralizadas aqui.
package config

import (
	"os"
	"strconv"
)

// Config representa todas as configurações da aplicação.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig contém as configurações do servidor HTTP.
type ServerConfig struct {
	Port string
	Mode string // "debug", "release", "test"
}

// DatabaseConfig contém as configurações de conexão com o SQL Server.
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// JWTConfig contém as configurações de autenticação JWT.
type JWTConfig struct {
	SecretKey          string
	AccessTokenExpiry  int // em horas
	RefreshTokenExpiry int // em dias
}

// Load carrega todas as configurações a partir de variáveis de ambiente.
// Valores padrão são usados quando as variáveis não estão definidas.
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "1433"),
			User:     getEnv("DB_USER", "sa"),
			Password: getEnv("DB_PASSWORD", "YourStrong@Passw0rd"),
			Database: getEnv("DB_NAME", "amigos_terceira_idade"),
		},
		JWT: JWTConfig{
			SecretKey:          getEnv("JWT_SECRET", "sua-chave-secreta-aqui-mude-em-producao"),
			AccessTokenExpiry:  getEnvAsInt("JWT_ACCESS_EXPIRY_HOURS", 24),
			RefreshTokenExpiry: getEnvAsInt("JWT_REFRESH_EXPIRY_DAYS", 7),
		},
	}
}

// getEnv retorna o valor da variável de ambiente ou o valor padrão.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt retorna o valor da variável de ambiente como inteiro ou o valor padrão.
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
