package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
	"github.com/vvjke314/mkc-backend/internal/pkg/filehandler"
)

// GetProjects godoc
// @Summary      Gets all customer projects
// @Description  Gets all customer projects
// @Tags         project
// @Produce      json
// @Security 	 BearerAuth
// @Success      200 {object} []ds.Project
// @Failure 500 {object} errorResponse
// @Router      /projects [get]
func (a *Application) GetProjects(c *gin.Context) {
	// Получаем JWT Токен
	tokenString := getJWT(c)
	// Парсим токен и получаем id клиента
	customerId, err := getJWTClaims(tokenString)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, "Can't parse JWT token")
		a.Log(err.Error())
		return
	}

	// Возвращеаем все проекты в ответ на запрос
	projects, err := a.repo.GetProjects(customerId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't handle request")
		a.Log(err.Error())
		return
	}

	c.JSON(http.StatusOK, projects)
}

// CreateProject godoc
// @Summary      Creates customer project
// @Description  Creates customer project in storage
// @Tags         project
// @Produce      json
// @Security 	 BearerAuth
// @Param data body ds.CreateProjectReq true "New project"
// @Success      200 {object} []ds.Project
// @Failure 500 {object} errorResponse
// @Router      /project [post]
func (a *Application) CreateProject(c *gin.Context) {
	req := &ds.CreateProjectReq{}
	// Анмаршалим тело запроса
	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Can't decode body params")
		a.Log(err.Error())
		return
	}

	// Получаем JWT Токен
	tokenString := getJWT(c)
	// Парсим токен и получаем id клиента
	customerId, err := getJWTClaims(tokenString)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, "Can't parse JWT token")
		a.Log(err.Error())
		return
	}

	if err = a.repo.GetProjectbyName(customerId, req.Name, &ds.Project{}); err == nil {
		newErrorResponse(c, http.StatusBadRequest, "Such project name already exists")
		return
	}

	// Подготавливаем структуру project для записи в БД
	project := ds.Project{
		Id:           uuid.New(),
		OwnerId:      uuid.MustParse(customerId),
		Capacity:     0,
		Name:         req.Name,
		CreationDate: time.Now(),
	}

	// Создаем папку проекта в нашем хранилище
	err = filehandler.CreateDir(project.Id.String())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't handle request")
		a.Log(err.Error())
		return
	}

	// Создаем запись о проекте в БД
	err = a.repo.CreateProject(project)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't handle request")
		a.Log(err.Error())
		return
	}

	// Возвращеаем все проекты в ответ на запрос
	projects, err := a.repo.GetProjects(project.OwnerId.String())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't handle request")
		a.Log(err.Error())
		return
	}

	c.JSON(http.StatusOK, projects)
}

// UpdateProjectName godoc
// @Summary      Updates project name
// @Description  Updates project name
// @Tags         project
// @Produce      json
// @Security 	 BearerAuth
// @Param project_id path string true "Project ID"
// @Param data body ds.UpdateProjectNameReq true "New project name"
// @Success      200 {object} []ds.Project
// @Failure 500 {object} errorResponse
// @Failure 401 {obejct} errorResponse
// @Failure 403 {object} errorResponse
// @Router      /project/{project_id} [put]
func (a *Application) UpdateProjectName(c *gin.Context) {
	req := &ds.UpdateProjectNameReq{}
	// Анмаршалим тело запроса
	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Can't decode body params")
		a.Log(err.Error())
		return
	}

	// Получаем JWT Токен
	tokenString := getJWT(c)
	// Парсим токен и получаем id клиента
	customerId, err := getJWTClaims(tokenString)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, "Can't parse JWT token")
		a.Log(err.Error())
		return
	}

	// // Проверка на доступ к работе с проектом, добавить промежуточное ПО
	// b, err := a.repo.AccessControl(customerId, c.Param("project_id"))
	// if !b || err != nil {
	// 	newErrorResponse(c, http.StatusForbidden, "You don't have permission to work with that project")
	// 	a.Log(err.Error())
	// 	return
	// }

	// Проверка на уже существующее имя проекта
	if err = a.repo.GetProjectbyName(customerId, req.Name, &ds.Project{}); err == nil {
		newErrorResponse(c, http.StatusBadRequest, "Such project already exists")
		return
	}

	// Обновляем имя проекта
	err = a.repo.UpdateProjectName(c.Param("project_id"), req.Name)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Can't update project name")
		a.Log(err.Error())
		return
	}

	// Возвращеаем все проекты в ответ на запрос
	projects, err := a.repo.GetProjects(customerId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't handle request")
		a.Log(err.Error())
		return
	}

	c.JSON(http.StatusOK, projects)
}

// AddParticipant godoc
// @Summary      Adds participant to project
// @Description  Adds participant to project
// @Tags         participants
// @Produce      json
// @Security 	 BearerAuth
// @Param project_id path string true "Project ID"
// @Param data body ds.AddParticipantReq true "CHANGE IT"
// @Success      200 {object} []ds.Customer
// @Failure 500 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Router      /participants/{project_id} [post]
func (a *Application) AddParticipant(c *gin.Context) {
	req := &ds.AddParticipantReq{}
	// Анмаршалим тело запроса
	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Can't decode body params")
		a.Log(err.Error())
		return
	}

	// Получаем JWT Токен
	tokenString := getJWT(c)
	// Парсим токен и получаем id клиента
	customerId, err := getJWTClaims(tokenString)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, "Can't parse JWT token")
		a.Log(err.Error())
		return
	}

	// Проверка на доступ к работе с проектом, добавить промежуточное ПО
	// b, err := a.repo.AccessControl(customerId, c.Param("project_id"))
	// if !b || err != nil {
	// 	newErrorResponse(c, http.StatusForbidden, "You don't have permission to work with that project")
	// 	a.Log(err.Error())
	// 	return
	// }

	// Получение пользователя и проверка на существования пользователя
	customer := &ds.Customer{}
	if err = a.repo.GetCustomerByEmail(req.ParticipantEmail, customer); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "No customer with such email")
		a.Log(err.Error())
		return
	}

	if customer.Id.String() == customerId {
		newErrorResponse(c, http.StatusBadRequest, "You owner of this project")
		a.Log("You owner of this project")
		return
	}

	// Хэш карта для хранения значения поля доступа в int типе
	customer_access := map[string]int{"просмотр": 0, "полный": 1}
	// Готовим строку для добавления в БД
	pa := ds.ProjectAccess{
		Id:             uuid.New(),
		ProjectId:      uuid.MustParse(c.Param("project_id")),
		CustomerId:     customer.Id,
		CustomerAccess: customer_access[req.CustomerAccess],
	}

	// Проверить не добавлен ли уже этот участник в проект
	err = a.repo.CheckParticipant(customer.Id.String(), c.Param("project_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		err = fmt.Errorf("[CheckParticipant] %w", err)
		a.Log(err.Error())
		return
	}

	// Добавляем участника в проект
	err = a.repo.CreateParticipant(pa)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "No such customer. Check your data")
		err = fmt.Errorf("[CreateParticipant] %w", err)
		a.Log(err.Error())
		return
	}

	// Получаем всех участников проекта
	customers, err := a.repo.GetParticipants(c.Param("project_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't get all project participants")
		err = fmt.Errorf("[GetParticipant] %w", err)
		a.Log(err.Error())
		return
	}

	c.JSON(http.StatusOK, customers)
}

// UpdateParticipantAccess godoc
// @Summary      Updates participant access in project
// @Description  Updates participant access in project
// @Tags         participants
// @Produce      json
// @Security 	 BearerAuth
// @Param project_id path string true "Project ID"
// @Param data body ds.UpdateParticipantAccessReq true "CHANGE IT"
// @Success      200 {object} []ds.Customer
// @Failure 500 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Router      /participants/{project_id} [put]
func (a *Application) UpdateParticipantAccess(c *gin.Context) {
	req := &ds.UpdateParticipantAccessReq{}
	// Анмаршалим тело запроса
	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Can't decode body params")
		a.Log(err.Error())
		return
	}

	// Хэш карта для хранения значения поля доступа в int типе
	customer_access := map[string]int{"просмотр": 0, "полный": 1}

	// Получение пользователя и проверка на существования пользователя
	customer := &ds.Customer{}
	if err = a.repo.GetCustomerByEmail(req.ParticipantEmail, customer); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "No customer with such email")
		a.Log(err.Error())
		return
	}

	// Проверить не добавлен ли уже этот участник в проект
	err = a.repo.CheckParticipant(customer.Id.String(), c.Param("project_id"))
	if err != nil {
		err = a.repo.UpdateParticipantAccess(customer.Id.String(), customer_access[req.CustomerAccess])
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, "Can't update customer access")
			err = fmt.Errorf("[UpdateParticipantAccess] %w", err)
			a.Log(err.Error())
			return
		}

		// Получаем всех участников проекта
		customers, err := a.repo.GetParticipants(c.Param("project_id"))
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, "Can't add participant")
			err = fmt.Errorf("[GetParticipant] %w", err)
			a.Log(err.Error())
			return
		}

		c.JSON(http.StatusOK, customers)
	}

	newErrorResponse(c, http.StatusBadRequest, err.Error())
	err = fmt.Errorf("[CheckParticipant] %w", err)
	a.Log(err.Error())
}

// DeleteParticipant godoc
// @Summary      Removes participant from project
// @Description  Removes participant from project
// @Tags         participants
// @Produce      json
// @Security 	 BearerAuth
// @Param project_id path string true "Project ID"
// @Param data body ds.DeleteParticipantReq true "CHANGE IT"
// @Success      200 {object} []ds.Customer
// @Failure 500 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Router      /participants/{project_id} [delete]
func (a *Application) DeleteParticipant(c *gin.Context) {
	req := &ds.DeleteParticipantReq{}
	// Анмаршалим тело запроса
	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Can't decode body params")
		a.Log(err.Error())
		return
	}

	// Получаем JWT Токен
	tokenString := getJWT(c)
	// Парсим токен и получаем id клиента
	customerId, err := getJWTClaims(tokenString)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, "Can't parse JWT token")
		a.Log(err.Error())
		return
	}

	// Получение пользователя и проверка на существования пользователя
	customer := &ds.Customer{}
	if err = a.repo.GetCustomerByEmail(req.ParticipantEmail, customer); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "No customer with such email")
		a.Log(err.Error())
		return
	}

	// Проверка не ввел ли клиент свой email
	if customer.Id.String() == customerId {
		newErrorResponse(c, http.StatusBadRequest, "You can't remove yourself project")
		a.Log("You can't remove yourself project")
		return
	}

	// Удаляем участника из проекта
	if err = a.repo.DeleteParticipant(customer.Id.String(), c.Param("project_id")); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't delete participant")
		err = fmt.Errorf("[DeleteParticipant] %w", err)
		a.Log(err.Error())
	}

	// Получаем всех участников проекта
	customers, err := a.repo.GetParticipants(c.Param("project_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Can't add participant")
		err = fmt.Errorf("[GetParticipant] %w", err)
		a.Log(err.Error())
		return
	}
	c.JSON(http.StatusOK, customers)
}

// DeleteProject godoc
// @Summary      Deletes project
// @Description  Deletes project
// @Tags         project
// @Produce      json
// @Security 	 BearerAuth
// @Param project_id path string true "Project ID"
// @Param data body ds.DeleteProjectReq true "CHANGE IT"
// @Success      200 {object} []ds.Project
// @Failure 500 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Router      /project/{project_id} [delete]
func (a *Application) DeleteProject(c *gin.Context) {
	req := &ds.DeleteProjectReq{}
	// Анмаршалим тело запроса
	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Can't decode body params")
		a.Log(err.Error())
		return
	}

	// Получаем JWT Токен
	tokenString := getJWT(c)
	// Парсим токен и получаем id клиента
	customerId, err := getJWTClaims(tokenString)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, "Can't parse JWT token")
		a.Log(err.Error())
		return
	}

	err = a.repo.DeleteProject(c.Param("project_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't delete project")
		a.Log(err.Error())
		return
	}

	// Возвращеаем все проекты в ответ на запрос
	projects, err := a.repo.GetProjects(customerId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't handle request")
		a.Log(err.Error())
		return
	}

	c.JSON(http.StatusOK, projects)
}
