package app

import "github.com/gin-gonic/gin"

// AttachAdmin
// @Summary      Attachs admin to project
// @Description  Attachs admin to project
// @Tags         auth
// @Produce      json
// @Security 	 BearerAuth
// @Success      200 {object} []ds.Project
// @Failure 500 {object} errorResponse
// @Router      /admin/project [post]
func (a *Application) AttachAdmin(c *gin.Context) {

}
