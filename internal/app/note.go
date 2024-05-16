package app

import (
	"encoding/json"
	"fmt"
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
	customerId := c.GetString("customerId")
	req := &ds.CreateNoteReq{}
	// Анмаршалим тело запроса
	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "can't decode body params")
		err = fmt.Errorf("[CreateNote][json.NewDecoder]: %w", err)
		a.Log(err.Error(), customerId)
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

	// Проверка на существование проекта
	if err := a.repo.GetNoteByName(note.Title, projectId); err == nil {
		newErrorResponse(c, http.StatusBadRequest, "such note title already exists in project")
		a.Log("[CreateNote][repo.GetNoteByName]: such note title already exists in project", customerId)
		return
	}

	// Создаем запись о заметке в БД
	err = a.repo.CreateNote(note)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't handle request")
		err = fmt.Errorf("[CreateNote][repo.CreateNote]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Возвращаем все заметки в проекте в ответ на запрос
	notes, err := a.repo.GetNotes(projectId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't handle request")
		err = fmt.Errorf("[CreateNote][repo.GetNotes]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	a.SuccessLog("[CreateNote]", customerId)
	c.JSON(http.StatusOK, notes)
}
