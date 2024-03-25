package service

import "github.com/gin-gonic/gin"

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

// CORSMiddleware
// мидлвейр для настройки политики CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") //localhost:3000
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Run
// Запускаем сервис на c помощью gin
func (s *Service) Run() {
	r := gin.Default()

	r.Use(CORSMiddleware())
}
