package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
)

// CreateNote godoc
// @Summary      Создание заявки в проекте
// @Description  Создание заявки в проекте
// @Tags         note
// @Produce      json
// @Security 	 BearerAuth
// @Param data body ds.CreateNoteReq true "New project"
// @Success      200 {object} []ds.Note
// @Failure 500 {object} errorResponse
// @Router      /project/:project_id/note [post]
func (a *Application) CreateNote(c *gin.Context) {
	req := &ds.CreateNoteReq{}
	// Анмаршалим тело запроса
	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Can't decode body params")
		a.Log(err.Error())
		return
	}

	projectId := c.GetString("projectId")

	// Подготавливаем структуру note для записи в БД
	note := ds.Note{
		Id:             uuid.New(),
		ProjectId:      uuid.MustParse(projectId),
		Title:          req.Title,
		Content:        req.Content,
		UpdateDatetime: time.Now(),
		Deadline:       req.Deadline,
		Overdue:        0,
	}

	// Создаем запись о заметке в БД
	err = a.repo.CreateNote(note)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't handle request")
		a.Log(err.Error())
		return
	}

	// Возвращаем все заметки в проекте в ответ на запрос
	notes, err := a.repo.GetNotes(projectId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't handle request")
		a.Log(err.Error())
		return
	}

	c.JSON(http.StatusOK, notes)
}
