package app

import "github.com/gin-gonic/gin"

// AttachAdmin
// @Summary      Прикрепляет администратора к проекту
// @Description  Прикрепляет администратора к выбраному проекту
// @Tags         auth
// @Produce      json
// @Security 	 BearerAuth
// @Success      200 {object} []ds.Project
// @Failure 500 {object} errorResponse
// @Router      /admin/project/{project_id} [post]
func (a *Application) AttachAdmin(c *gin.Context) {

}

// GetAllUnattachedProjects
// @Summary      Все проекты которые еще не прикреплены
// @Description  Возвращает все проекты которые еще не прикреплены
// @Tags         auth
// @Produce      json
// @Security 	 BearerAuth
// @Success      200 {object} []ds.Project
// @Failure 500 {object} errorResponse
// @Router      /admin/unattached [get]
func (a *Application) GetAllUnattachedProjects(c *gin.Context) {

}

// GetAllAttachedProjects
// @Summary      Все проекты которые прикреплены к администратору
// @Description  Возравщает все проекты которые прикреплены к администратору
// @Tags         auth
// @Produce      json
// @Security 	 BearerAuth
// @Success      200 {object} []ds.Project
// @Failure 500 {object} errorResponse
// @Router      /admin/attached [get]
func (a *Application) GetAllAttachedProjects(c *gin.Context) {

}

// SumbitEmail
// @Summary      Отправляет сообщение пользователю на почту
// @Description  Отправляет сообщение пользователю на почту
// @Tags         auth
// @Produce      json
// @Security 	 BearerAuth
// @Success      200 {object} []ds.Project
// @Failure 500 {object} errorResponse
// @Router      /admin/:project_id/send [post]
func (a *Application) SumbitEmail(c *gin.Context) {

}
