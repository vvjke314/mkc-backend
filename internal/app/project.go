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
// @Summary      Возвращаем все проекты пользователя
// @Description  Возращает все проекты пользователя
// @Tags         project
// @Produce      json
// @Security 	 BearerAuth
// @Success      200 {object} []ds.Project
// @Failure 500 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Router      /projects [get]
func (a *Application) GetProjects(c *gin.Context) {
	customerId := c.GetString("customer_id")

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
// @Summary      Создает проект пользователю
// @Description  Создает проект пользователю
// @Tags         project
// @Produce      json
// @Security 	 BearerAuth
// @Param data body ds.CreateProjectReq true "New project"
// @Success      200 {object} []ds.Project
// @Failure 500 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 403 {object} errorResponse
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

	// Проверяем проект на существование
	if err = a.repo.GetProjectbyName(customerId, req.Name, &ds.Project{}); err == nil {
		newErrorResponse(c, http.StatusBadRequest, "such project name already exists")
		a.Log("[CreateProject]:such project already exists", customerId)
		return
	}

	// Получаем все проекты
	prjcts, err := a.repo.GetProjects(customerId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "database can't exec your query")
		a.Log("[CreateProject]:such project already exists", customerId)
		return
	}

	isSub := c.GetBool("isSubscription")
	// Проверка на подписку
	if !isSub && len(prjcts) >= 5 {
		newErrorResponse(c, http.StatusForbidden, "you already in 5 projects you need to upgrade your account status")
		a.Log("[CreateProject]:you already in 5 projects you need to upgrade your account status", customerId)
		return
	}

	// Подготавливаем структуру project для записи в БД
	project := ds.Project{
		Id:           uuid.New(),
		OwnerId:      uuid.MustParse(customerId),
		Capacity:     52428800,
		Name:         req.Name,
		CreationDate: time.Now(),
	}

	// Для пользователей с платной подпиской больше места для хранения данных
	if isSub {
		project.Capacity = 56000000
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
// @Summary      Обновляет имя проекта
// @Description  Обновляет имя проекта
// @Tags         project
// @Produce      json
// @Security 	 BearerAuth
// @Param project_id path string true "Project name"
// @Param data body ds.UpdateProjectNameReq true "New project name"
// @Success      200 {object} []ds.Project
// @Failure 500 {object} errorResponse
// @Failure 401 {object} errorResponse
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

// DeleteProject godoc
// @Summary      Удаляет проект
// @Description  Удаляет проект
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
