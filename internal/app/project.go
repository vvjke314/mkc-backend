package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
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
		newErrorResponse(c, http.StatusBadRequest, "Can't parse JWT token")
	}
	// Подготавливаем структуру project для записи в БД
	project := ds.Project{
		Id:           uuid.New(),
		OwnerId:      uuid.MustParse(customerId),
		Capacity:     0,
		Name:         req.Name,
		CreationDate: time.Now(),
	}

	err = a.repo.CreateProject(project)
	if err != nil {
		// Обработать ошибку
		return
	}

	projects, err := a.repo.GetProjects(project.OwnerId.String())
	if err != nil {
		// Обработать ошибку
		return
	}

	c.JSON(http.StatusOK, projects)
}

// UpdateProjectName godoc
// @Summary      Updates project name
// @Description  Updates project name
// @Tags         auth
// @Produce      json
// @Security 	 BearerAuth
// @Param project_id path string true "Project ID"
// @Param data body ds.UpdateProjectNameReq true "New project name"
// @Success      200 {object} []ds.Project
// @Failure 500 {object} errorResponse
// @Router      /project/{project_id} [put]
func (a *Application) UpdateProjectName(c *gin.Context) {

}

// AddParticipant godoc
// @Summary      Adds participant to project
// @Description  Adds participant to project
// @Tags         auth
// @Produce      json
// @Security 	 BearerAuth
// @Param project_id path string true "Project ID"
// @Param data body ds.AddParticipantReq true "CHANGE IT"
// @Success      200 {object} []ds.Customer
// @Failure 500 {object} errorResponse
// @Router      /participants/{project_id} [post]
func (a *Application) AddParticipant(c *gin.Context) {

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
