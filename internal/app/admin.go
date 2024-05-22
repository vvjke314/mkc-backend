package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vvjke314/mkc-backend/internal/pkg/crypt"
	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
)

// SignUpAdmin
// @Summary      Добавляет администратора на сервис
// @Description  Добавляет администратора на сервис
// @Tags         administrator
// @Produce      json
// @Param data body ds.SignUpAdmin true "Информация о администраторе"
// @Success      200 {object} successResponse
// @Failure 500 {object} errorResponse
// @Failure 400 {object} errorResponse
// @Router      /admin/signup [post]
func (a *Application) SignUpAdmin(c *gin.Context) {
	req := &ds.SignUpAdmin{}

	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Can't decode body params")
		return
	}

	if req.Password == "" {
		newErrorResponse(c, http.StatusBadRequest, "Password is empty")
		return
	}

	if req.Name == "" {
		newErrorResponse(c, http.StatusBadRequest, "Name is empty")
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

	admin := ds.Administrator{
		Id:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: password,
	}
	err = a.repo.SignUpAdministrator(admin)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed with signing up")
		return
	}

	// Успешный ответ на запрос
	newSuccessResponse(c, 200, "successfully signed up")
}

// AttachAdmin
// @Summary      Прикрепляет администратора к проекту
// @Description  Прикрепляет администратора к выбраному проекту
// @Tags         administrator
// @Produce      json
// @Security 	 BasicAuth
// @Param project_id path string true "Уникальный идентификатор проекта"
// @Success      200 {object} successResponse
// @Failure 500 {object} errorResponse
// @Router      /admin/project/{project_id} [get]
func (a *Application) AttachAdmin(c *gin.Context) {
	projectId := c.Param("project_id")
	adminId := c.GetString("adminId")

	err := a.repo.SetAdministrator(adminId, projectId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed with signing up")
		err = fmt.Errorf("[AttachAdmin][repo.SetAdministrator]: %w", err)
		a.Log(err.Error(), adminId)
		return
	}

	// Успешный ответ на запрос
	a.Log("[AttachAdmin]", adminId)
	newSuccessResponse(c, 200, "successfully attached")
}

// GetAllUnattachedProjects
// @Summary      Все проекты которые еще не прикреплены
// @Description  Возвращает все проекты которые еще не прикреплены
// @Tags         administrator
// @Produce      json
// @Security 	 BasicAuth
// @Success      200 {object} []ds.Project
// @Failure 500 {object} errorResponse
// @Router      /admin/unattached [get]
func (a *Application) GetAllUnattachedProjects(c *gin.Context) {
	adminId := c.GetString("adminId")

	projects, err := a.repo.GetAllUnattachedProjects()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't get all untachhed projects")
		err = fmt.Errorf("[GetAllUnattachedProjects][repo.GetAllUnattachedProjects]: %w", err)
		a.Log(err.Error(), adminId)
		return
	}

	// Успешный ответ на запрос
	a.Log("[GetAllUnattachedProjects]", adminId)
	c.JSON(200, projects)
}

// GetAllAttachedProjects
// @Summary      Все проекты которые прикреплены к администратору
// @Description  Возвращает все проекты которые прикреплены к администратору
// @Tags         administrator
// @Produce      json
// @Security 	 BasicAuth
// @Success      200 {object} []ds.Project
// @Failure 500 {object} errorResponse
// @Router      /admin/attached [get]
func (a *Application) GetAllAttachedProjects(c *gin.Context) {
	adminId := c.GetString("adminId")

	projects, err := a.repo.GetAllAttachedProjects(adminId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't get all atached projects")
		err = fmt.Errorf("[GetAllAttachedProjects][repo.GetAllAttachedProjects]: %w", err)
		a.Log(err.Error(), adminId)
		return
	}

	// Успешный ответ на запрос
	a.Log("[GetAllAttachedProjects]", adminId)
	c.JSON(200, projects)
}

// GetCustomerEmail
// @Summary      Получает электронную почту пользователя, владеющего проектом
// @Description  Получает электронную почту пользователя, владеющего проектом
// @Tags         administrator
// @Produce      json
// @Security 	 BasicAuth
// @Param project_id path string true "Уникальный идентификатор проекта"
// @Success      200 {object} ds.GetCustomerEmailResponse
// @Failure 500 {object} errorResponse
// @Router      /admin/{project_id}/send [post]
func (a *Application) GetCustomerEmail(c *gin.Context) {
	projectId := c.Param("project_id")
	adminId := c.GetString("adminId")

	email, err := a.repo.GetCustomerEmail(projectId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't get customer email")
		err = fmt.Errorf("[GetCustomerEmail][repo.GetCustomerEmail]: %w", err)
		a.Log(err.Error(), adminId)
		return
	}

	// Возвращаем почту
	a.Log("[GetCustomerEmail]", adminId)
	c.JSON(200, ds.GetCustomerEmailResponse{
		Email: email,
	})
}
