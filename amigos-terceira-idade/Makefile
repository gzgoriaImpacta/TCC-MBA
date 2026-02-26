# Makefile para comandos comuns do projeto
# Use: make <comando>

.PHONY: help run build test docker-up docker-down clean

# Variáveis
APP_NAME=amigos-terceira-idade
MAIN_PATH=./cmd/api

# Comando padrão - mostra ajuda
help:
	@echo "Comandos disponíveis:"
	@echo "  make run          - Executa a aplicação localmente"
	@echo "  make build        - Compila a aplicação"
	@echo "  make test         - Executa os testes"
	@echo "  make test-cover   - Executa testes com cobertura"
	@echo "  make docker-up    - Sobe os containers (SQL Server + API)"
	@echo "  make docker-db    - Sobe apenas o SQL Server"
	@echo "  make docker-down  - Para os containers"
	@echo "  make clean        - Limpa arquivos de build"
	@echo "  make deps         - Baixa as dependências"
	@echo "  make lint         - Executa o linter"

# Baixa as dependências
deps:
	go mod download
	go mod tidy

# Executa a aplicação localmente
run:
	go run $(MAIN_PATH)/main.go

# Compila a aplicação
build:
	go build -o bin/$(APP_NAME) $(MAIN_PATH)/main.go

# Executa os testes
test:
	go test -v ./...

# Executa testes com cobertura
test-cover:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Relatório de cobertura gerado: coverage.html"

# Sobe todos os containers
docker-up:
	docker-compose up -d

# Sobe apenas o SQL Server (para desenvolvimento local)
docker-db:
	docker-compose up -d sqlserver

# Para os containers
docker-down:
	docker-compose down

# Para e remove volumes
docker-clean:
	docker-compose down -v

# Limpa arquivos de build
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Executa o linter (requer golangci-lint instalado)
lint:
	golangci-lint run ./...

# Formata o código
fmt:
	go fmt ./...
	goimports -w .

# Cria o banco de dados (execute após docker-db)
create-db:
	docker exec -i amigos-sqlserver //opt/mssql-tools18/bin/sqlcmd \
		-S localhost -U sa -P 'YourStrong@Passw0rd' \
		-C -Q "IF NOT EXISTS (SELECT name FROM sys.databases WHERE name = 'amigos_terceira_idade') CREATE DATABASE amigos_terceira_idade"