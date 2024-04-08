package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
	"github.com/vvjke314/mkc-backend/internal/pkg/filehandler"
)

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
		return
	}

	// Получаем JWT Токен
	tokenString := getJWT(c)
	// Парсим токен и получаем id клиента
	customerId, err := getJWTClaims(tokenString)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, "Can't parse JWT token")
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
		a.logger.Error().Msg(err.Error())
		return
	}

	// Создаем запись о проекте в БД
	err = a.repo.CreateProject(project)
	if err != nil {
		a.logger.Error().Msg(err.Error())
		return
	}

	// Возвращеаем все проекты в ответ на запрос
	projects, err := a.repo.GetProjects(project.OwnerId.String())
	if err != nil {
		a.logger.Error().Msg(err.Error())
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
		return
	}

	// Получаем JWT Токен
	tokenString := getJWT(c)
	// Парсим токен и получаем id клиента
	customerId, err := getJWTClaims(tokenString)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, "Can't parse JWT token")
		return
	}

	// Проверка на доступ к работе с проектом, добавить промежуточное ПО
	b, err := a.repo.AccessControl(customerId, c.Param("project_id"))
	if !b || err != nil {
		newErrorResponse(c, http.StatusForbidden, "You don't have permission to work with that project")
		return
	}

	// Проверка на уже существующее имя проекта
	if err = a.repo.GetProjectbyName(customerId, req.Name, &ds.Project{}); err == nil {
		newErrorResponse(c, http.StatusBadRequest, "Such project already exists")
		return
	}

	// Обновляем имя проекта
	err = a.repo.UpdateProjectName(c.Param("project_id"), req.Name)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Can't update project name")
		return
	}

	// Возвращеаем все проекты в ответ на запрос
	projects, err := a.repo.GetProjects(customerId)
	if err != nil {
		a.logger.Error().Msg(err.Error())
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
		return
	}

	// Получаем JWT Токен
	tokenString := getJWT(c)
	// Парсим токен и получаем id клиента
	customerId, err := getJWTClaims(tokenString)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, "Can't parse JWT token")
		return
	}

	// Проверка на доступ к работе с проектом, добавить промежуточное ПО
	b, err := a.repo.AccessControl(customerId, c.Param("project_id"))
	if !b || err != nil {
		newErrorResponse(c, http.StatusForbidden, "You don't have permission to work with that project")
		return
	}

	// Получение пользователя и проверка на существования пользователя
	customer := &ds.Customer{}
	if err = a.repo.GetCustomerByEmail(req.ParticipantLogin, customer); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "No customer with such email")
		return
	}

	pa := ds.ProjectAccess{
		Id:             uuid.New(),
		ProjectId:      uuid.MustParse(c.Param("project_id")),
		CustomerId:     customer.Id,
		CustomerAccess: 0,
	}

	// Добавляем участника в проект
	// Добавить участника проект
	//
	err = a.repo.CreateParticipant(pa)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Can't add participant")
		return
	}

	customers, err := a.repo.GetParticipants(c.Param("project_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Can't add participant")
		return
	}

	c.JSON(http.StatusOK, customers)
}

// DeleteParticipant godoc
// @Summary      Removes participant from project
// @Description  Removes participant from project
// @Tags         auth
// @Produce      json
// @Security 	 BearerAuth
// @Param project_id path string true "Project ID"
// @Param data body ds.DeleteParticipantReq true "CHANGE IT"
// @Success      200 {object} []ds.Customer
// @Failure 500 {object} errorResponse
// @Router      /participants/{project_id} [delete]
func (a *Application) DeleteParticipant(c *gin.Context) {

}

// DeleteProject godoc
// @Summary      Deletes project
// @Description  Deletes project
// @Tags         auth
// @Produce      json
// @Security 	 BearerAuth
// @Param project_id path string true "Project ID"
// @Param data body ds.DeleteProjectReq true "CHANGE IT"
// @Success      200 {object} []ds.Customer
// @Failure 500 {object} errorResponse
// @Router      /project/{project_id} [delete]
func (a *Application) DeleteProject(c *gin.Context) {

}
