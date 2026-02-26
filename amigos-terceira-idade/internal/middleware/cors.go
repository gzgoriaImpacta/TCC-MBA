// Package middleware contém os middlewares HTTP da aplicação.
package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware configura os headers CORS para permitir requisições de outros domínios.
// Necessário para que o frontend (React Native Web, React) consiga acessar a API.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Permite requisições de qualquer origem em desenvolvimento
		// Em produção, substitua "*" pelos domínios específicos
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Platform, X-App-Version")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		// Responde imediatamente para requisições OPTIONS (preflight)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
