package app

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// BasicAuthMiddleware промежуточное ПО для проверки является ли отправитель администратором
func (a *Application) BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем заголовок авторизации
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Header("WWW-Authenticate", `Basic realm="Authorization Required"`)
			newErrorResponse(c, http.StatusUnauthorized, "you must authorized to do this action")
			a.Log("must authorized", "BasicAuthMiddleware")
			return
		}

		// Проверяем формат авторизации
		if !strings.HasPrefix(authHeader, "Basic ") {
			newErrorResponse(c, http.StatusUnauthorized, "you must authorized to do this action")
			a.Log("must authorized", "BasicAuthMiddleware")
			return
		}

		// Декодируем пароль
		decoded, err := base64.StdEncoding.DecodeString(authHeader[6:])
		if err != nil {
			newErrorResponse(c, http.StatusUnauthorized, "you must authorized to do this action")
			a.Log("must authorized", "BasicAuthMiddleware")
			return
		}

		// Разделяем пароль и имя администратора
		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 {
			newErrorResponse(c, http.StatusUnauthorized, "you must authorized to do this action")
			a.Log("must authorized", "BasicAuthMiddleware")
			return
		}

		// Получаем валидные данные
		hashedPassword, err := a.repo.GetValidCredentials(creds[0])
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, "can't get valid data")
			a.Log(fmt.Errorf("[GetValidCredentials]: can't get valid data %w", err).Error(), "BasicAuthMiddleware")
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds[1])); err != nil {
			newErrorResponse(c, http.StatusUnauthorized, "bad authorize password")
			a.Log(fmt.Errorf("[crypt.HashPassword]: bad authorize password %w", err).Error(), "BasicAuthMiddleware")
			return
		}

		adminId, err := a.repo.GetAdminId(creds[0], hashedPassword)
		if err != nil {
			newErrorResponse(c, http.StatusUnauthorized, "can't get administrator ID")
			a.Log("[GetAdminId] can't get administrator ID", "BasicAuthMiddleware")
			return
		}

		c.Set("adminId", adminId)

		c.Next()
	}
}

// FullAccessControl промежуточное ПО для проверки доступа к работе с проектом
func (a *Application) FullAccessControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем JWT Токен
		tokenString := getJWT(c)
		// Парсим токен и получаем id клиента
		customerId, err := getJWTClaims(tokenString)
		if err != nil {
			newErrorResponse(c, http.StatusForbidden, "can't parse JWT token")
			a.Log(err.Error(), customerId)
			return
		}
		projectId := c.Param("project_id")
		c.Set("customerId", customerId)
		c.Set("projectId", projectId)

		b, err := a.repo.AccessControl(customerId, projectId, 1)
		if !b {
			if err == nil {
				newErrorResponse(c, http.StatusForbidden, "you don't have permission to work with that project")
				a.Log("customer don't have permission to work with project", customerId)
				return
			} else {
				newErrorResponse(c, http.StatusInternalServerError, "database can't parse you query")
				a.Log(err.Error(), customerId)
				return
			}
		}
		c.Next()
	}
}

// AccessControl промежуточное ПО для проверки доступа к работе с проектом
func (a *Application) AccessControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем JWT Токен
		tokenString := getJWT(c)
		// Парсим токен и получаем id клиента
		customerId, err := getJWTClaims(tokenString)
		if err != nil {
			newErrorResponse(c, http.StatusForbidden, "can't parse JWT token")
			return
		}
		projectId := c.Param("project_id")
		c.Set("customerId", customerId)
		c.Set("projectId", projectId)

		b, err := a.repo.AccessControl(customerId, projectId, 0)
		if !b {
			if err == nil {
				newErrorResponse(c, http.StatusForbidden, "you don't have permission to work with that project")
				a.Log("customer don't have permission to work with project", customerId)
				return
			} else {
				newErrorResponse(c, http.StatusInternalServerError, "database can't parse you query")
				a.Log(err.Error(), customerId)
				return
			}
		}
		c.Next()
	}
}

// CORSMiddleware промежуточное ПО для настройки политики CORS
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

// AuthMiddleware промежуточное ПО для проверки на наличие JWT-токена
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем JWT-токен из заголовка Authorization
		tokenRawString := c.GetHeader("Authorization")
		if tokenRawString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		// Удбираем Bearer из заголовка аутентификации
		tokenSplitString := strings.Split(tokenRawString, " ")
		tokenString := tokenSplitString[1]

		// Проверяем, есть ли токен в списке активных токенов
		if _, ok := activeTokens[tokenString]; !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Продолжаем выполнение запроса
		c.Next()
	}
}

// Получаем наш токен в виде строки
func getJWT(c *gin.Context) string {
	tokenRawString := c.GetHeader("Authorization")
	// Удбираем Bearer из заголовка аутентификации
	tokenSplitString := strings.Split(tokenRawString, " ")
	tokenString := tokenSplitString[1]
	return tokenString
}

// Парсим и проверяем токен и получаем payload из него
func getJWTClaims(tokenString string) (string, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи токена
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Возвращаем секретный ключ для проверки подписи токена
		return []byte("mkc-forever-alone"), nil
	})
	if err != nil {
		// Обработка ошибки парсинга токена
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		// Получаем значение поля "id" из токена
		id := claims["id"].(string)
		return id, nil
	}
	return "", fmt.Errorf("invalid token")

}

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, errMessage string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{errMessage})
}

type successResponse struct {
	Message string `json:"message"`
}

func newSuccessResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, successResponse{Message: message})
}
