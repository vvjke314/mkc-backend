package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	customerId := c.GetString("customerId")

	// Возвращаем все проекты в ответ на запрос
	projects, err := a.repo.GetProjects(customerId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't handle request")
		err = fmt.Errorf("[GetProjects][repo.GetProjects]:%w", err)
		a.Log(err.Error(), customerId)
		return
	}

	a.SuccessLog("[GetProjects]", customerId)
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
	// Получаем JWT Токен
	tokenString := getJWT(c)
	// Парсим токен и получаем id клиента
	customerId, err := getJWTClaims(tokenString)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, "can't parse JWT token")
		a.Log(err.Error(), customerId)
		return
	}
	req := &ds.CreateProjectReq{}
	// Анмаршалим тело запроса
	err = json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "can't decode body params")
		err = fmt.Errorf("[CreateProject][json.NewDecoder]:%w", err)
		a.Log(err.Error(), customerId)
		return
	}

	if err = a.repo.GetProjectbyName(customerId, req.Name, &ds.Project{}); err == nil {
		newErrorResponse(c, http.StatusBadRequest, "such project name already exists")
		a.Log("[CreateProject]:such project already exists", customerId)
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
		newErrorResponse(c, http.StatusInternalServerError, "can't handle request")
		err = fmt.Errorf("[CreateProject][filehandler.CreateDir]:%w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Создаем запись о проекте в БД
	err = a.repo.CreateProject(project)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't handle request")
		err = fmt.Errorf("[CreateProject][repo.CreateProject]:%w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Возвращеаем все проекты в ответ на запрос
	projects, err := a.repo.GetProjects(project.OwnerId.String())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't handle request")
		err = fmt.Errorf("[CreateProject][repo.GetProjects]:%w", err)
		a.Log(err.Error(), customerId)
		return
	}

	a.SuccessLog("[CreateProject]", customerId)
	c.JSON(http.StatusOK, projects)
}

// UpdateProjectName godoc
// @Summary      Updates project name
// @Description  Updates project name
// @Tags         project
// @Produce      json
// @Security 	 BearerAuth
// @Param project_id path string true "Project name"
// @Param data body ds.UpdateProjectNameReq true "New project name"
// @Success      200 {object} []ds.Project
// @Failure 500 {object} errorResponse
// @Failure 401 {obejct} errorResponse
// @Failure 403 {object} errorResponse
// @Router      /project/{project_id} [put]
func (a *Application) UpdateProjectName(c *gin.Context) {
	customerId := c.GetString("customerId")
	projectId := c.GetString("projectId")

	req := &ds.UpdateProjectNameReq{}
	// Анмаршалим тело запроса
	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "can't decode body params")
		err = fmt.Errorf("[UpdateProjectName][json.NewDecoder]:%w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Проверка на уже существующее имя проекта
	if err = a.repo.GetProjectbyName(customerId, req.Name, &ds.Project{}); err == nil {
		newErrorResponse(c, http.StatusBadRequest, "such project already exists")
		err = fmt.Errorf("[UpdateProjectName][repo.GetProjectByName]:such project alredy exist")
		a.Log(err.Error(), customerId)
		return
	}

	// Обновляем имя проекта
	err = a.repo.UpdateProjectName(projectId, req.Name)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "can't update project name")
		err = fmt.Errorf("[UpdateProjectName][repo.UpdateProjectName]:%w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Возвращеаем все проекты в ответ на запрос
	projects, err := a.repo.GetProjects(customerId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't handle request")
		err = fmt.Errorf("[UpdateProjectName][repo.GetProjects]:%w", err)
		a.Log(err.Error(), customerId)
		return
	}

	a.SuccessLog("[UpdateProjectName]", customerId)
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
	customerId := c.GetString("customerId")
	req := &ds.AddParticipantReq{}
	// Анмаршалим тело запроса
	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "can't decode body params")
		err = fmt.Errorf("[AddParticipant][json.NewDecoder]:%w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Получение пользователя и проверка на существования пользователя
	customer := &ds.Customer{}
	if err = a.repo.GetCustomerByEmail(req.ParticipantEmail, customer); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "no customer with such email")
		err = fmt.Errorf("[AddParticipant][repo.GetCustomerByEmail]:%w", err)
		a.Log(err.Error(), customerId)
		return
	}

	if customer.Id.String() == customerId {
		newErrorResponse(c, http.StatusBadRequest, "you owner of this project")
		err = fmt.Errorf("[AddParticipant]:you owner of this project")
		a.Log(err.Error(), customerId)
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
		err = fmt.Errorf("[AddParticipant][repo.CheckParticipants]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Добавляем участника в проект
	err = a.repo.CreateParticipant(pa)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "no such customer. Check your data")
		err = fmt.Errorf("[AddParticipant][repo.CreateParticipant]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Получаем всех участников проекта
	customers, err := a.repo.GetParticipants(c.Param("project_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't get all project participants")
		err = fmt.Errorf("[AddParticipant][repo.GetParticipant]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	a.SuccessLog("[AddParticipant]", customerId)
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
	customerId := c.GetString("customerId")
	req := &ds.UpdateParticipantAccessReq{}
	// Анмаршалим тело запроса
	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "can't decode body params")
		err = fmt.Errorf("[UpdateParticipantAccess][json.NewDecoder]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Хэш карта для хранения значения поля доступа в int типе
	customer_access := map[string]int{"просмотр": 0, "полный": 1}

	// Получение пользователя и проверка на существования пользователя
	customer := &ds.Customer{}
	if err = a.repo.GetCustomerByEmail(req.ParticipantEmail, customer); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "no customer with such email")
		err = fmt.Errorf("[UpdateParticipantAccess][repo.GetCustomerByEmail]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Проверить не добавлен ли уже этот участник в проект
	err = a.repo.CheckParticipant(customer.Id.String(), c.Param("project_id"))
	if err != nil {
		err = a.repo.UpdateParticipantAccess(customer.Id.String(), customer_access[req.CustomerAccess])
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, "Can't update customer access")
			err = fmt.Errorf("[UpdateParticipantAccess][repo.UpdateParticipantAccess]: %w", err)
			a.Log(err.Error(), customerId)
			return
		}

		// Получаем всех участников проекта
		customers, err := a.repo.GetParticipants(c.Param("project_id"))
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, "Can't add participant")
			err = fmt.Errorf("[UpdateParticipantAccess][GetParticipant]: %w", err)
			a.Log(err.Error(), customerId)
			return
		}
		a.SuccessLog("[UpdateParticipantAccess]", customerId)
		c.JSON(http.StatusOK, customers)
	}

	newErrorResponse(c, http.StatusBadRequest, err.Error())
	err = fmt.Errorf("[UpdateParticipantAccess][CheckParticipant]: %w", err)
	a.Log(err.Error(), customerId)
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
	customerId := c.GetString("customerId")

	req := &ds.DeleteParticipantReq{}
	// Анмаршалим тело запроса
	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "can't decode body params")
		err = fmt.Errorf("[DeleteParticipant][json.NewDecoder]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Получение пользователя и проверка на существования пользователя
	customer := &ds.Customer{}
	if err = a.repo.GetCustomerByEmail(req.ParticipantEmail, customer); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "no customer with such email")
		err = fmt.Errorf("[DeleteParticipant][repo.GetCustomerByEmail]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Проверка не ввел ли клиент свой email
	if customer.Id.String() == customerId {
		newErrorResponse(c, http.StatusBadRequest, "you can't remove yourself project")
		a.Log("[DeleteParticipant]: you can't remove yourself project", customerId)
		return
	}

	// Удаляем участника из проекта
	if err = a.repo.DeleteParticipant(customer.Id.String(), c.Param("project_id")); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't delete participant")
		err = fmt.Errorf("[DeleteParticipant][repo.DeleteParticipant]: %w", err)
		a.Log(err.Error(), customerId)
	}

	// Получаем всех участников проекта
	customers, err := a.repo.GetParticipants(c.Param("project_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Can't add participant")
		err = fmt.Errorf("[DeleteParticipant][repo.GetParticipant]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	a.SuccessLog("[DeleteParticipant]", customerId)
	c.JSON(http.StatusOK, customers)
}

// DeleteProject godoc
// @Summary      Deletes project
// @Description  Deletes project
// @Tags         project
// @Produce      json
// @Security 	 BearerAuth
// @Param project_id path string true "Project ID"
// @Success      200 {object} []ds.Project
// @Failure 500 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Router      /project/{project_id} [delete]
func (a *Application) DeleteProject(c *gin.Context) {
	customerId := c.GetString("customerId")
	projectId := c.GetString("projectId")

	err := a.repo.DeleteProject(projectId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't delete project")
		err = fmt.Errorf("[DeleteProject][repo.DeleteProject]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Возвращеаем все проекты в ответ на запрос
	projects, err := a.repo.GetProjects(customerId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't handle request")
		err = fmt.Errorf("[DeleteProject][repo.GetProjects]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Удаляем файл из локальной директории
	os.RemoveAll(filehandler.Path + projectId)

	a.SuccessLog("[DeleteProject]", customerId)
	c.JSON(http.StatusOK, projects)
}

// GetProjectInfo получить информацию о проекте
// @Summary      Получаем информацию о содержании проекта
// @Description  Получаем массив всех файлов и заметок
// @Tags         project
// @Produce      json
// @Security 	 BearerAuth
// @Param project_id path string true "Идентификатор проекта"
// @Success      200 {object} ds.ProjectData
// @Failure 500 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Router      /project/{project_id} [get]
func (a *Application) GetProjectInfo(c *gin.Context) {
	projectId := c.Param("project_id")
	customerId := c.GetString("customerId")

	files, err := a.repo.GetFiles(projectId)
	if err != nil && err != pgx.ErrNoRows {
		newErrorResponse(c, http.StatusInternalServerError, "can't handle request")
		err = fmt.Errorf("[DeleteProject][repo.GetProjects]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	notes, err := a.repo.GetNotes(projectId)
	if err != nil && err != pgx.ErrNoRows {
		newErrorResponse(c, http.StatusInternalServerError, "can't handle request")
		err = fmt.Errorf("[DeleteProject][repo.GetProjects]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	projectData := ds.ProjectData{
		Notes: notes,
		Files: files,
	}

	a.SuccessLog("[GetProjectInfo]", customerId)
	c.JSON(200, projectData)
}
