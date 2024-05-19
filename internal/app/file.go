package app

import (
	"encoding/json"
	"fmt"
	"io"
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

// UploadFile godoc
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
func (a *Application) UploadFile(c *gin.Context) {
	// Получаем проект ID из запроса
	projectId := c.GetString("projectId")
	customerId := c.GetString("customerId")

	// Получаем файл из запроса
	file, err := c.FormFile("file")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "no file uploaded")
		err = fmt.Errorf("[UploadFile][gin.Context.FormFile]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Генерируем уникальное имя файла
	filename := file.Filename
	fileNames := strings.Split(filename, ".")
	// Проверяем файл на существование
	if err := a.repo.CheckFileExistence(fileNames[0], "."+fileNames[1], projectId); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		err = fmt.Errorf("[UploadFile][repo.CheckFileExistence]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Сохраняем файл на сервере
	if err := c.SaveUploadedFile(file, filehandler.Path+projectId+"/"+filename); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "no file uploaded")
		err = fmt.Errorf("[UploadFile][gin.Context.SaveUploadedFile]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Добавляем запись о файле в базу данных
	newFile := ds.File{
		Id:             uuid.New(),
		ProjectId:      uuid.MustParse(projectId),
		Filename:       fileNames[0],
		Extension:      filepath.Ext(file.Filename),
		Size:           int(int(file.Size) / 1000), //размер в КБ
		FilePath:       filehandler.Path + projectId + "/" + filename,
		UpdateDatetime: time.Now(),
	}

	// Добавляем информацию о файле
	err = a.repo.CreateFile(newFile)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "no file added")
		err = fmt.Errorf("[UploadFile][repo.CreateFile]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Получаем все файлы из проекта
	files, err := a.repo.GetFiles(projectId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't get files")
		err = fmt.Errorf("[UploadFile][repo.GetFiles]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	a.SuccessLog("[UploadFile]", customerId)
	c.JSON(http.StatusOK, files)
}

// UploadFiles godoc
// @Summary Загрузить файлы
// @Description Загружает файлы на сервер
// @Tags file
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param project_id path string true "Идентификатор проекта"
// @Param files formData file true "Файлы для загрузки"
// @Success 200 {object} []ds.File
// @Failure 500 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Router /project/{project_id}/files [post]
func (a *Application) UploadFiles(c *gin.Context) {
	projectId := c.GetString("projectId")
	customerId := c.GetString("customerId")

	// Парсим форму с несколькими файлами
	err := c.Request.ParseMultipartForm(5 << 22) // Размер до 5 GB
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "error parsing form")
		err = fmt.Errorf("[UploadFiles][gin.Context.Request.ParseMultipartForm]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Получаем все файлы из формы
	files := c.Request.MultipartForm.File["files"]

	// Проходим по каждому файлу
	for _, file := range files {
		// Генерируем уникальное имя файла
		filename := file.Filename
		fileNames := strings.Split(filename, ".")

		// Сохраняем файл на сервере
		if err := c.SaveUploadedFile(file, filehandler.Path+projectId+"/"+filename); err != nil {
			newErrorResponse(c, http.StatusInternalServerError, "error saving file")
			err = fmt.Errorf("[UploadFiles][gin.Context.SaveUploadedFile]: %w", err)
			a.Log(err.Error(), customerId)
			return
		}

		// Добавляем запись о файле в базу данных
		newFile := ds.File{
			Id:             uuid.New(),
			ProjectId:      uuid.MustParse(projectId),
			Filename:       fileNames[0],
			Extension:      filepath.Ext(file.Filename),
			Size:           int(int(file.Size) / 1000), //размер в КБ
			FilePath:       filehandler.Path + projectId + "/" + filename,
			UpdateDatetime: time.Now(),
		}

		// Добавляем информацию о файле
		err = a.repo.CreateFile(newFile)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, "error adding file to database")
			err = fmt.Errorf("[UploadFiles][repo.CreateFile]: %w", err)
			a.Log(err.Error(), customerId)
			return
		}
	}

	// Получаем все файлы из проекта
	fs, err := a.repo.GetFiles(projectId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't get files")
		err = fmt.Errorf("[UploadFiles][repo.GetFiles]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	a.SuccessLog("[UploadFile]", customerId)
	c.JSON(http.StatusOK, fs)
}

// DeleteFile godoc
// @Summary Удалить файл
// @Description Удаляет файл с сервера и из БД
// @Tags file
// @Security 	 BearerAuth
// @Produce json
// @Param project_id path string true "Идентификатор проекта"
// @Param  data body ds.DeleteFileReq true "Структура хранящая тело запроса для удаления файла"
// @Failure 500 {object} errorResponse
// @Failure 401 {obejct} errorResponse
// @Failure 403 {object} errorResponse
// @Router /project/{project_id}/file [delete]
func (a *Application) DeleteFile(c *gin.Context) {
	projectId := c.GetString("projectId")
	customerId := c.GetString("customerId")
	req := &ds.DeleteFileReq{}
	// Анмаршалим тело запроса
	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "can't decode body params")
		err = fmt.Errorf("[DeleteFile][json.NewDecoder]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Получаем искомый файл
	file := &ds.File{}
	err = a.repo.GetFileByName(req.Filename, req.Extension, projectId, file)
	if err != nil {
		if err == pgx.ErrNoRows {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			err = fmt.Errorf("[DeleteFile][repo.GetFileByName]: %w", err)
			a.Log(err.Error(), customerId)
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		err = fmt.Errorf("[DeleteFile][repo.GetFileByName]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Удаляем файл из БД
	err = a.repo.DeleteFile(file.Id.String())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't delete file")
		err = fmt.Errorf("[DeleteFile][repo.GetFileByName]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Удаляем файл из хранилища
	err = os.Remove(file.FilePath)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't remove file from storage")
		err = fmt.Errorf("[DeleteFile][os.Remove]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Получаем массив из оставшихся файлов в проекте
	files, err := a.repo.GetFiles(file.ProjectId.String())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't scan all files from project")
		err = fmt.Errorf("[DeleteFile][repo.GetFiles]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	a.SuccessLog("[DeleteFile]", customerId)
	c.JSON(http.StatusOK, files)
}

// DownloadFile обрабатывает запрос на скачивание файла с сервера
// @Summary Скачать файл
// @Description Скачивает файл с сервера
// @Tags file
// @Accept json
// @Security BearerAuth
// @Produce octet-stream
// @Param project_id path string true "Идентификатор проекта"
// @Param file_id path string true "Идентификатор файла"
// @Success 200 {file} octet-stream
// @Failure 500 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Router /project/{project_id}/file/{file_id} [get]
func (a *Application) DownloadFile(c *gin.Context) {
	customerId := c.GetString("customerId")
	projectId := c.GetString("projectId")
	fileId := c.Param("file_id")
	file := &ds.File{}
	err := a.repo.GetFileById(fileId, file)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "file not found")
		err = fmt.Errorf("[DownloadFile][repo.GetFileById]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Формируем путь к файлу на сервере
	filePath := filehandler.Path + projectId + "/" + file.Filename + file.Extension

	// Проверяем существование файла
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		// Если файл не найден, возвращаем ошибку
		newErrorResponse(c, http.StatusNotFound, "file not found")
		err = fmt.Errorf("[DownloadFile][os.IsNotExist]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Открываем файл для чтения
	f, err := os.Open(filePath)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "file not found")
		err = fmt.Errorf("[DownloadFile][os.Open]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}
	defer f.Close()

	// Определяем MIME-тип по расширению файла
	contentType := "application/octet-stream" // MIME-тип по умолчанию

	extension := file.Extension
	if extension != "" {
		// Откидываем точку и переводим в нижний регистр
		extension = strings.ToLower(extension[1:])
		switch extension {
		case "pdf":
			contentType = "application/pdf"
		case "txt":
			contentType = "text/plain"
		case "jpg", "jpeg":
			contentType = "image/jpeg"
		case "png":
			contentType = "image/png"
		}
	}

	c.Header("Content-Description", "File Transfer")
	// Устанавливаем заголовок Content-Disposition для указания имени файла при скачивании
	c.Header("Content-Disposition", "attachment; filename="+file.Filename)
	// MIME-тип для бинарных данных
	c.Header("Content-Type", contentType)

	// Копируем содержимое файла в ответ HTTP
	if _, err := io.Copy(c.Writer, f); err != nil {
		newErrorResponse(c, http.StatusNotFound, "failed to copy file content to response")
		err = fmt.Errorf("[DownloadFile][io.Copy]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}
}

// GetFiles отдает все файлы из проекта
// @Summary Просмотреть все файлы проекта
// @Description Получить все файлы проекта
// @Tags file
// @Security 	 BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param project_id path string true "Идентификатор проекта"
// @Success 200 {object} []ds.File
// @Failure 500 {object} errorResponse
// @Failure 401 {obejct} errorResponse
// @Failure 403 {object} errorResponse
// @Router /project/{project_id}/files [get]
func (a *Application) GetFiles(c *gin.Context) {
	projectId := c.Param("project_id")
	customerId := c.GetString("customerId")
	// Получаем все файлы из проекта
	files, err := a.repo.GetFiles(projectId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't get files")
		err = fmt.Errorf("[GetFiles][repo.GetFiles]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	a.SuccessLog("GetFiles", customerId)
	c.JSON(http.StatusOK, files)
}
