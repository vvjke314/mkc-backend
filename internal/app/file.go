package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
	"github.com/vvjke314/mkc-backend/internal/pkg/filehandler"
)

// CreateFile загружает файл на сервер
// @Summary Загрузить файл
// @Description Загружает файл на сервер
// @Tags file
// @Security 	 BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param project_id path string true "Идентификатор проекта"
// @Param file formData file true "Файл для загрузки"
// @Success 200 {object} []ds.File
// @Failure 500 {object} errorResponse
// @Failure 401 {obejct} errorResponse
// @Failure 403 {object} errorResponse
// @Router /project/{project_id}/file [post]
func (a *Application) CreateFile(c *gin.Context) {
	// Получаем проект ID из запроса
	projectId := c.Param("project_id")

	// Получаем файл из запроса
	file, err := c.FormFile("file")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "No file uploaded")
		return
	}

	// Генерируем уникальное имя файла
	filename := file.Filename
	fileNames := strings.Split(filename, ".")
	// Проверяем файл на существование
	if err := a.repo.CheckFileExistence(fileNames[0], "."+fileNames[1], projectId); err != nil {
		err = fmt.Errorf("[repo.CheckFileExistence] %w", err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		a.Log(err.Error())
		return
	}

	// Сохраняем файл на сервере
	if err := c.SaveUploadedFile(file, filehandler.Path+projectId+"/"+filename); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "No file uploaded")
		return
	}

	// Добавляем запись о файле в базу данных
	newFile := ds.File{
		Id:             uuid.New(),
		ProjectId:      uuid.MustParse(projectId),
		Filename:       fileNames[0],
		Extension:      filepath.Ext(file.Filename),
		Size:           int(int(file.Size) / 1000),
		FilePath:       filehandler.Path + projectId + "/" + filename,
		UpdateDatetime: time.Now(),
	}

	// Добавляем информацию о файле
	err = a.repo.CreateFile(newFile)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "No file added")
		err = fmt.Errorf("[repo.CreateFile] %w", err)
		a.Log(err.Error())
		return
	}

	// Получаем все файлы из проекта
	files, err := a.repo.GetFiles(projectId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't get files")
		err = fmt.Errorf("[repo.GetFiles] %w", err)
		a.Log(err.Error())
		return
	}

	c.JSON(http.StatusOK, files)
}

// DeleteFile удаляет файл с сервера
// @Summary Удалить файл
// @Description Удаляет файл с сервера
// @Tags file
// @Security 	 BearerAuth
// @Produce json
// @Param project_id path string true "Идентификатор проекта"
// @Param  data body ds.DeleteFileReq true "CHANGE IT"
// @Failure 500 {object} errorResponse
// @Failure 401 {obejct} errorResponse
// @Failure 403 {object} errorResponse
// @Router /project/{project_id}/file [delete]
func (a *Application) DeleteFile(c *gin.Context) {
	projectId := c.Param("project_id")
	req := &ds.DeleteFileReq{}
	// Анмаршалим тело запроса
	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Can't decode body params")
		a.Log(err.Error())
		return
	}

	// Получаем искомый файл
	file := &ds.File{}
	err = a.repo.GetFileByName(req.Filename, req.Extension, projectId, file)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = fmt.Errorf("[repo.GetFileByName] %w", err)
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			a.Log(err.Error())
			return
		}
		err = fmt.Errorf("[repo.GetFileByName] %w", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		a.Log(err.Error())
		return
	}

	// Удаляем файл из БД
	err = a.repo.DeleteFile(file.Id.String())
	if err != nil {
		err = fmt.Errorf("[repo.DeleteFile] %w", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		a.Log(err.Error())
		return
	}

	// Удаляем файл из хранилища
	err = os.Remove(file.FilePath)
	if err != nil {
		err = fmt.Errorf("[os.Remove] %w", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		a.Log(err.Error())
		return
	}

	// Получаем массив из оставшихся файлов в проекте
	files, err := a.repo.GetFiles(file.ProjectId.String())
	if err != nil {
		err = fmt.Errorf("[repo.GetFiles] %w", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		a.Log(err.Error())
		return
	}

	c.JSON(http.StatusOK, files)
}

// UpdateFileName меняет имя файла на сервере
// @Summary Изменяет имя файла
// @Description Изменяет имя файла на сервере
// @Tags file
// @Security 	 BearerAuth
// @Produce json
// @Param project_id path string true "Идентификатор проекта"
// @Param  data body ds.UpdateFileNameReq true "CHANGE IT"
// @Failure 500 {object} errorResponse
// @Failure 401 {obejct} errorResponse
// @Failure 403 {object} errorResponse
// @Router /project/{project_id} [post]
func (a *Application) UpdateFileName(c *gin.Context) {

}

// GetFile обрабатывает запрос на скачивание файла с сервера
// @Summary Скачать файл
// @Description Скачивает файл с сервера
// @Tags File
// @Accept json
// @Produce octet-stream
// @Param project_id path string true "Идентификатор проекта"
// @Param file_name path string true "Имя файла"
// @Success 200 {file} octet-stream
// @Failure 500 {object} errorResponse
// @Failure 401 {obejct} errorResponse
// @Failure 403 {object} errorResponse
// @Router /project/{project_id}/file/{file_name} [get]
func (a *Application) GetFile(c *gin.Context) {

}
