package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vvjke314/mkc-backend/internal/pkg/crypt"
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
func createToken(login, id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"login": login,
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	})

	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", fmt.Errorf("[jwt.Token.SignedString] %w", err)
	}
	return tokenString, nil
}

// Login godoc
// @Summary      Логинит пользователя
// @Description  Логинит пользователя
// @Tags         auth
// @Produce      json
// @Param data body ds.LoginCustomerReq true "Customer data"
// @Success      200 {object} AuthToken
// @Failure 500 {object} errorResponse
// @Failure 400 {object} errorResponse
// @Router      /login [post]
func (a *Application) Login(c *gin.Context) {
	var creds ds.LoginCustomerReq
	if err := c.ShouldBindJSON(&creds); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	customer := ds.Customer{}
	if err := a.repo.GetCustomerByCredentials(creds, &customer); err != nil {
		if err == pgx.ErrNoRows {
			newErrorResponse(c, http.StatusBadRequest, "no such customer")
			err = fmt.Errorf("[Login][repo.GetCustomerByCredentials]: %w", err)
			a.Log(err.Error(), "[Authentification]")
			return
		}
		newErrorResponse(c, http.StatusBadRequest, "invalid password")
		err = fmt.Errorf("[Login][repo.GetCustomerByCredentials]: %w", err)
		a.Log(err.Error(), "[Authentification]")
		return
	}

	token, err := createToken(customer.Login, customer.Id.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	activeTokens[token] = true

	// Возврат JWT-токена клиенту
	a.SuccessLog("successfuly signed up", customer.Id.String())
	c.JSON(http.StatusOK, AuthToken{Token: token})
}

// SignUp godoc
// @Summary      Регистрирует пользователя
// @Description  Регистрирует пользователя
// @Tags         auth
// @Produce      json
// @Param data body ds.SignUpCustomerReq true "Customer data"
// @Success      200 {object} AuthToken
// @Failure 500 {object} errorResponse
// @Failure 400 {object} errorResponse
// @Router      /signup [post]
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
		return
	}

	if req.Email == "" {
		newErrorResponse(c, http.StatusBadRequest, "Email is empty")
		return
	}

	password, err := crypt.HashPassword(req.Password)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Bad password entered")
		return
	}

	customer := ds.Customer{
		Id:         uuid.New(),
		FirstName:  req.FirstName,
		SecondName: req.SecondName,
		Login:      req.Login,
		Email:      req.Email,
		Password:   password,
		Type:       0,
	}
	err = a.repo.SignUpCustomer(customer)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Failed with signing up. Customer with entered data alredy exist")
		return
	}
	// Выполняем регистрацию пользователя и, в случае успеха, создать JWT-токен
	token, err := createToken(customer.Login, customer.Id.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}
	activeTokens[token] = true

	// Возврат JWT-токена клиенту
	c.JSON(http.StatusOK, AuthToken{Token: token})
}

// Logout godoc
// @Summary      Разлогинивает пользователя
// @Description  Разлогинивает пользователя
// @Tags         auth
// @Produce      json
// @Security 	 BearerAuth
// @Success      200 {object} successResponse
// @Failure 403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router      /logout [get]
func (a *Application) Logout(c *gin.Context) {
	// Получаем JWT-токен из заголовка авторизации
	tokenRawString := c.GetHeader("Authorization")
	if tokenRawString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	// Удбираем Bearer из заголовка аутентификации
	tokenSplitString := strings.Split(tokenRawString, " ")
	tokenString := tokenSplitString[1]

	// Удаляем токен из списка активных токенов
	delete(activeTokens, tokenString)

	// Отдаем сообщение об успешном выходе.
	newSuccessResponse(c, http.StatusOK, "Logout successful")
}
