// Package middleware contém os middlewares HTTP da aplicação.
// Middlewares são funções que processam a requisição antes de chegar ao handler.
package middleware

import (
	"net/http"
	"strings"
	"amigos-terceira-idade/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware cria um middleware de autenticação JWT.
// Valida o token no header Authorization e adiciona os dados do usuário ao contexto.
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtém o header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "MISSING_TOKEN",
					"message": "Token de autenticação não fornecido",
				},
			})
			c.Abort()
			return
		}

		// Verifica o formato do header (Bearer token)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_TOKEN_FORMAT",
					"message": "Formato do token inválido. Use: Bearer {token}",
				},
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Valida o token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_TOKEN",
					"message": "Token inválido ou expirado",
				},
			})
			c.Abort()
			return
		}

		// Adiciona os dados do usuário ao contexto para uso nos handlers
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_type", claims.UserType)

		c.Next()
	}
}
