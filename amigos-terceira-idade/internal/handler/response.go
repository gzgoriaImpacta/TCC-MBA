// Package handler contém os handlers HTTP da aplicação.
// Handlers são responsáveis por receber requisições e retornar respostas.
package handler

import (
	"github.com/gin-gonic/gin"
)

// Response representa a estrutura padrão de resposta da API.
// Todas as respostas seguem este formato para consistência.
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *MetaInfo   `json:"meta,omitempty"`
}

// ErrorInfo contém detalhes sobre um erro.
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// MetaInfo contém informações adicionais como paginação.
type MetaInfo struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

// SuccessResponse retorna uma resposta de sucesso padronizada.
func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

// SuccessResponseWithMeta retorna uma resposta de sucesso com metadados.
func SuccessResponseWithMeta(c *gin.Context, statusCode int, data interface{}, meta *MetaInfo) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// ErrorResponse retorna uma resposta de erro padronizada.
func ErrorResponse(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}
