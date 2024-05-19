package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
)

// CreateNote создает заметку
// @Summary      Создание заявки в проекте
// @Description  Создание заявки в проекте и добавление записи в БД
// @Tags         note
// @Produce      json
// @Security 	 BearerAuth
// @Param data body ds.CreateNoteReq true "New project"
// @Success      200 {object} []ds.Note
// @Failure 500 {object} errorResponse
// @Router      /project/{project_id}/note [post]
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

// DeleteNote удаляет заметку из проекта
// @Summary Удалить заметку
// @Description Удаляет заметки из проекта и БД
// @Tags note
// @Security 	 BearerAuth
// @Produce json
// @Param project_id path string true "Идентификатор проекта"
// @Param note_id path string true "Идентификатор заметки"
// @Success 200 {object} []ds.Note
// @Failure 500 {object} errorResponse
// @Failure 401 {obejct} errorResponse
// @Failure 403 {object} errorResponse
// @Router /project/{project_id}/note/{note_id} [delete]
func (a *Application) DeleteNote(c *gin.Context) {
	projectId := c.GetString("projectId")
	customerId := c.GetString("customerId")
	noteId := c.Param("note_id")

	n := &ds.Note{}
	// Получаем искомую заметку
	if err := a.repo.GetNoteById(noteId, n); err != nil {
		if err == pgx.ErrNoRows {
			newErrorResponse(c, http.StatusBadRequest, "no such note in this project")
			err = fmt.Errorf("[DeleteNote][repo.GetNoteById]: %w", err)
			a.Log(err.Error(), customerId)
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, "can't handle request")
		err = fmt.Errorf("[DeleteNote][repo.GetNoteById]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Удаляем заметку из БД
	if err := a.repo.DeleteNote(noteId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't handle request")
		err = fmt.Errorf("[DeleteNote][repo.DeleteNote]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Получаем массив из оставшихся заметок в проекте
	notes, err := a.repo.GetNotes(projectId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't handle request")
		err = fmt.Errorf("[DeleteNote][repo.GetNotes]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	a.SuccessLog("[DeleteNote]", customerId)
	c.JSON(http.StatusOK, notes)
}

// UpdateNoteDeadline обновляем дедлайн заметки
// @Summary Обновить дедлайн заметки
// @Description Обновляет дедлайн заметки в БД
// @Tags note
// @Security 	 BearerAuth
// @Produce json
// @Param project_id path string true "Идентификатор проекта"
// @Param note_id path string true "Идентификатор заметки"
// @Success 200 {object} []ds.Note
// @Failure 500 {object} errorResponse
// @Failure 400 {obejct} errorResponse
// @Router /project/{project_id}/note/{note_id} [put]
func (a *Application) UpdateNoteDeadline(c *gin.Context) {
	projectId := c.GetString("projectId")
	customerId := c.GetString("customerId")
	noteId := c.Param("note_id")

	req := &ds.UpdateNoteDeadlineReq{}
	// Анмаршалим тело запроса
	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "can't decode body params")
		err = fmt.Errorf("[UpdateNoteDeadline][json.NewDecoder]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	n := &ds.Note{}
	// Получаем искомую заметку
	if err := a.repo.GetNoteById(noteId, n); err != nil {
		if err == pgx.ErrNoRows {
			newErrorResponse(c, http.StatusBadRequest, "no such note in this project")
			err = fmt.Errorf("[UpdateNoteDeadline][repo.GetNoteById]: %w", err)
			a.Log(err.Error(), customerId)
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, "can't handle request")
		err = fmt.Errorf("[UpdateNoteDeadline][repo.GetNoteById]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Обновляем значение дедлайна заметки
	if err := a.repo.UpdateNoteDeadLine(noteId, req.Deadline); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't handle request")
		err = fmt.Errorf("[UpdateNoteDeadline][repo.UpdateNoteDeadLine]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Получаем массив из заметок в проекте
	notes, err := a.repo.GetNotes(projectId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't handle request")
		err = fmt.Errorf("[UpdateNoteDeadline][repo.GetNotes]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	a.SuccessLog("[UpdateNoteDeadline]", customerId)
	c.JSON(http.StatusOK, notes)
}
