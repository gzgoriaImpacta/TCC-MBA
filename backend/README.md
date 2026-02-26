# Amigos da Terceira Idade - API

> Conectando gerações com empatia e tecnologia

Plataforma digital que conecta **voluntários** a **idosos** e **instituições**, promovendo companhia, acolhimento e vínculos afetivos com segurança e praticidade.

## Stack Técnica

| Camada | Tecnologia |
|--------|------------|
| **Backend** | Go 1.21 + Gin |
| **Banco de Dados** | SQL Server |
| **Autenticação** | JWT (golang-jwt) |
| **ORM** | GORM |
| **Frontend (consumidor)** | React Native (iOS/Android) + React Web |

## Início Rápido

### Pré-requisitos

- Go 1.21+
- Docker e Docker Compose
- Make (opcional, mas recomendado)

### 1. Clone o repositório

```bash
git clone https://github.com/fmeireles/amigos-terceira-idade.git
cd amigos-terceira-idade
```

### 2. Configure as variáveis de ambiente

```bash
cp env.example .env
# Edite o arquivo .env conforme necessário
```

### 3. Suba o banco de dados

```bash
# Sobe apenas o SQL Server
make docker-db

# Aguarde ~30 segundos para o SQL Server iniciar
# Crie o banco de dados
make create-db
```

### 4. Execute a aplicação

```bash
# Baixa as dependências
make deps

# Executa a aplicação
make run
```

A API estará disponível em `http://localhost:8080`

### 5. Teste a API

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Lista interesses disponíveis
curl http://localhost:8080/api/v1/interests
```

## Endpoints da API

### Core MVP (Implementados)

#### Autenticação
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `POST` | `/api/v1/auth/register` | Cadastro de usuário |
| `POST` | `/api/v1/auth/login` | Login |
| `POST` | `/api/v1/auth/refresh` | Renovar tokens |

#### Usuários
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/users/me` | Meu perfil |
| `PUT` | `/api/v1/users/me` | Atualizar perfil |
| `DELETE` | `/api/v1/users/me` | Desativar conta |
| `GET` | `/api/v1/users/:id` | Ver perfil de usuário |

#### Interesses
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/interests` | Lista todos os interesses |
| `GET` | `/api/v1/interests/:id` | Busca interesse por ID |

#### Pareamento (Matching)
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/matching/suggestions` | Sugestões de pareamento |
| `POST` | `/api/v1/matching/connect` | Criar conexão |
| `GET` | `/api/v1/matching/connections` | Minhas conexões |
| `POST` | `/api/v1/matching/connections/:id/accept` | Aceitar conexão |
| `POST` | `/api/v1/matching/connections/:id/reject` | Rejeitar conexão |

#### Agendamentos
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `POST` | `/api/v1/appointments` | Criar agendamento |
| `GET` | `/api/v1/appointments` | Meus agendamentos |
| `GET` | `/api/v1/appointments/upcoming` | Próximos agendamentos |
| `GET` | `/api/v1/appointments/:id` | Detalhes do agendamento |
| `POST` | `/api/v1/appointments/:id/accept` | Aceitar convite |
| `POST` | `/api/v1/appointments/:id/decline` | Recusar convite |
| `DELETE` | `/api/v1/appointments/:id` | Cancelar agendamento |

#### Convites
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/invitations/received` | Convites recebidos |
| `GET` | `/api/v1/invitations/sent` | Convites enviados |

### Future Ready (Preparados para Evolução)

Os seguintes endpoints estão **planejados** e serão implementados conforme necessidade:

- [ ] `POST /api/v1/users/me/emergency-contact` - Contato de emergência
- [ ] `POST /api/v1/volunteers/:id/verification` - Verificação de voluntário
- [ ] `GET /api/v1/volunteers/:id/achievements` - Badges/conquistas
- [ ] `GET /api/v1/volunteers/:id/stats` - Estatísticas (horas dedicadas)
- [ ] `PUT /api/v1/users/me/availability` - Disponibilidade semanal
- [ ] `GET /api/v1/notifications` - Notificações
- [ ] `GET /api/v1/connections/:id/messages` - Chat
- [ ] `GET /api/v1/admin/dashboard` - Dashboard administrativo

## Exemplos de Uso

### Cadastro de Voluntário

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ricardo Silva",
    "email": "ricardo@email.com",
    "password": "senha123",
    "age": 35,
    "bio": "Gosto de conversar, ouvir histórias e fazer companhia.",
    "user_type": "VOLUNTEER",
    "interest_ids": ["uuid-interesse-1", "uuid-interesse-2"]
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "ricardo@email.com",
    "password": "senha123"
  }'
```

### Buscar Sugestões de Pareamento

```bash
curl http://localhost:8080/api/v1/matching/suggestions \
  -H "Authorization: Bearer SEU_TOKEN_AQUI"
```

### Criar Agendamento

```bash
curl -X POST http://localhost:8080/api/v1/appointments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN_AQUI" \
  -d '{
    "target_id": "uuid-do-idoso",
    "date": "2026-01-25T15:00:00Z",
    "duration_minutes": 30,
    "notes": "Primeira conversa"
  }'
```

## Estrutura do Projeto

```
amigos-terceira-idade/
├── cmd/
│   └── api/
│       └── main.go              # Ponto de entrada
├── internal/
│   ├── config/
│   │   └── config.go            # Configurações
│   ├── domain/
│   │   ├── user.go              # Entidade User
│   │   ├── interest.go          # Entidade Interest
│   │   ├── connection.go        # Entidade Connection
│   │   └── appointment.go       # Entidade Appointment
│   ├── repository/
│   │   ├── user_repository.go
│   │   ├── interest_repository.go
│   │   ├── connection_repository.go
│   │   └── appointment_repository.go
│   ├── service/
│   │   ├── auth_service.go      # Autenticação e JWT
│   │   ├── user_service.go
│   │   ├── interest_service.go
│   │   ├── matching_service.go  # Lógica de pareamento
│   │   └── appointment_service.go
│   ├── handler/
│   │   ├── router.go            # Configuração de rotas
│   │   ├── response.go          # Respostas padronizadas
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── interest_handler.go
│   │   ├── matching_handler.go
│   │   └── appointment_handler.go
│   └── middleware/
│       ├── auth.go              # Validação JWT
│       └── cors.go              # CORS para Web
├── pkg/
│   └── database/
│       └── sqlserver.go         # Conexão SQL Server
├── test/                        # Testes
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod
└── README.md
```

## Comandos Úteis

```bash
make help           # Lista todos os comandos
make run            # Executa localmente
make build          # Compila
make test           # Executa testes
make test-cover     # Testes com cobertura
make docker-up      # Sobe todos os containers
make docker-db      # Sobe apenas o SQL Server
make docker-down    # Para os containers
make lint           # Executa o linter
make fmt            # Formata o código
```

## Testes

```bash
# Executar todos os testes
make test

# Executar com cobertura
make test-cover

# Abrir relatório de cobertura
open coverage.html
```

## Roadmap

### Fase 1 - MVP ✅
- [x] Cadastro e Login
- [x] Perfil de usuário
- [x] Interesses
- [x] Pareamento básico
- [x] Agendamentos
- [x] Convites

### Fase 2 - Segurança (Planejado)
- [ ] Verificação de voluntário
- [ ] Contato de emergência
- [ ] Avaliação pós-conversa
- [ ] Botão de pânico

### Fase 3 - Engajamento (Planejado)
- [ ] Notificações push
- [ ] Badges/gamificação
- [ ] Chat simples
- [ ] Streak de voluntariado

### Fase 4 - Escala (Planejado)
- [ ] Dashboard admin
- [ ] Integração Google Meet
- [ ] Integração WhatsApp
- [ ] Disponibilidade recorrente

## Contribuindo

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -m 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## Licença

Este projeto foi desenvolvido como parte do MBA e está disponível para fins educacionais.

---

**"Todo idoso merece ser lembrado, escutado e abraçado — mesmo que a distância."**

Vamos juntos espalhar afeto com tecnologia!
