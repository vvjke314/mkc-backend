package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
)

var activeTokens = make(map[string]bool)

var (
	JwtSecret = []byte("mkc-forever-alone")
)

type AuthToken struct {
	Token string `json:"token"`
}

// createToken создание JWT-токена
func createToken(login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": login,
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	})

	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", fmt.Errorf("[jwt.Token.SignedString] %w", err)
	}
	return tokenString, nil
}

type Credentials struct {
	Login    string
	Password string
}

func Login(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Здесь следует выполнить проверку учетных данных пользователя и, в случае успеха, создать JWT-токен
	token, err := createToken(creds.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	activeTokens[token] = true

	// Возврат JWT-токена клиенту
	c.JSON(http.StatusOK, AuthToken{Token: token})
}

// Signup создание пользователя
func (a *Application) Signup(c *gin.Context) {
	req := &ds.SignUpCustomerReq{}

	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Can't decode body params")
		return
	}

	if req.Password == "" {
		newErrorResponse(c, http.StatusBadRequest, "Password is empty")
		return
	}

	if req.FirstName == "" {
		newErrorResponse(c, http.StatusBadRequest, "Firstname is empty")
		return
	}

	if req.SecondName == "" {
		newErrorResponse(c, http.StatusBadRequest, "Secondname is empty")
		return
	}

	if req.Login == "" {
		newErrorResponse(c, http.StatusBadRequest, "Login is empty")
	}

	if req.Email == "" {
		newErrorResponse(c, http.StatusBadRequest, "Email is empty")
		return
	}

	// Здесь следует выполнить регистрацию пользователя и, в случае успеха, создать JWT-токен
	token, err := createToken(req.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	activeTokens[token] = true

	err = a.repo.SignUpCustomer(ds.Customer{
		Id:         uuid.New(),
		FirstName:  req.FirstName,
		SecondName: req.SecondName,
		Login:      req.Login,
		Email:      req.Email,
		Password:   req.Password,
		Type:       0,
	})

	// Возврат JWT-токена клиенту
	c.JSON(http.StatusOK, AuthToken{Token: token})
}

// authMiddleware промежуточное ПО для проверки на наличие JWT-токена
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем JWT-токен из заголовка Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

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

func Logout(c *gin.Context) {
	// Получаем JWT-токен из заголовка авторизации
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Удаляем токен из списка активных токенов
	delete(activeTokens, tokenString)

	// В данном примере просто возвращаем сообщение об успешном выходе.
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
