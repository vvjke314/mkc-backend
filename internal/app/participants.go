package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
)

// AddParticipant godoc
// @Summary      Добавляет участника в проект
// @Description  Добавляет участника в проект
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
		if err == pgx.ErrNoRows {
			newErrorResponse(c, http.StatusBadRequest, "no customer with such email")
			err = fmt.Errorf("[AddParticipant][repo.GetCustomerByEmail]:%w", err)
			a.Log(err.Error(), customerId)
		} else {
			newErrorResponse(c, http.StatusInternalServerError, "can't parse your query")
			err = fmt.Errorf("[AddParticipant][repo.GetCustomerByEmail]:%w", err)
			a.Log(err.Error(), customerId)
		}
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
// @Summary      Обновить доступ участнику проекта
// @Description  Обновить доступ участнику проекта. В поле CustomerAccess вводить либо "полный" либо "просмотр"
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
			if err == pgx.ErrNoRows {
				newErrorResponse(c, http.StatusBadRequest, "no customer with such email")
				err = fmt.Errorf("[AddParticipant][repo.GetCustomerByEmail]:%w", err)
				a.Log(err.Error(), customerId)
			} else {
				newErrorResponse(c, http.StatusInternalServerError, "can't parse your query")
				err = fmt.Errorf("[AddParticipant][repo.GetCustomerByEmail]:%w", err)
				a.Log(err.Error(), customerId)
			}
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
		return
	}

	newErrorResponse(c, http.StatusBadRequest, "no such customer in project check your data")
	err = fmt.Errorf("[UpdateParticipantAccess][CheckParticipant]: %w", err)
	a.Log(err.Error(), customerId)
}

// DeleteParticipant godoc
// @Summary      Убрать участника из проекта
// @Description  Убрать участника из проекта
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
	projectId := c.Param("project_id")

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
		if err == pgx.ErrNoRows {
			newErrorResponse(c, http.StatusBadRequest, "no customer with such email")
			err = fmt.Errorf("[DeleteParticipant][repo.GetCustomerByEmail]:%w", err)
			a.Log(err.Error(), customerId)
			return
		} else {
			newErrorResponse(c, http.StatusInternalServerError, "can't parse your query")
			err = fmt.Errorf("[DeleteParticipant][repo.GetCustomerByEmail]:%w", err)
			a.Log(err.Error(), customerId)
			return
		}
	}

	// Проверка на существование пользователя в проекте
	if err := a.repo.CheckParticipant(customer.Id.String(), projectId); err == nil {
		newErrorResponse(c, http.StatusBadRequest, "no customer in project with such email")
		err = fmt.Errorf("[DeleteParticipant][repo.CheckParticipant]:%w", err)
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
		return
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

// GetAllParticipants godoc
// @Summary      показать всех участников проекта
// @Description  показать всех участников проекта, включая его создателя
// @Tags         participants
// @Produce      json
// @Security 	 BearerAuth
// @Param project_id path string true "Project ID"
// @Success      200 {object} []ds.Customer
// @Failure 500 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Router      /participants/{project_id} [get]
func (a *Application) GetAllParticipants(c *gin.Context) {
	customerId := c.GetString("customerId")

	// Получаем всех участников проекта
	customers, err := a.repo.GetParticipants(c.Param("project_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Can't add participant")
		err = fmt.Errorf("[GetAllParticipants][repo.GetParticipant]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	a.SuccessLog("[GetAllParticipants]", customerId)
	c.JSON(http.StatusOK, customers)
}
