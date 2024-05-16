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

// UploadFile загружает файл на сервер
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
		Size:           int(int(file.Size) / 1000), //размер в КБ
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

// UploadFiles загружает несколько файлов на сервер
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

	// Парсим форму с несколькими файлами
	err := c.Request.ParseMultipartForm(5 << 22) // Размер до 5 GB
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Error parsing form")
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
			newErrorResponse(c, http.StatusInternalServerError, "Error saving file")
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
			newErrorResponse(c, http.StatusInternalServerError, "Error adding file to database")
			return
		}
	}

	// Получаем все файлы из проекта
	fs, err := a.repo.GetFiles(projectId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't get files")
		return
	}

	c.JSON(http.StatusOK, fs)
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
	projectId := c.Param("project_id")
	fileId := c.Param("file_id")
	file := &ds.File{}
	err := a.repo.GetFileById(fileId, file)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Формируем путь к файлу на сервере
	filePath := filehandler.Path + projectId + "/" + file.Filename

	// Проверяем существование файла
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		// Если файл не найден, возвращаем ошибку
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Открываем файл для чтения
	f, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy file content to response"})
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
	// Получаем все файлы из проекта
	files, err := a.repo.GetFiles(projectId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't get files")
		err = fmt.Errorf("[repo.GetFiles] %w", err)
		a.Log(err.Error())
		return
	}

	// Успешное завершение запроса
	c.JSON(http.StatusOK, files)
}
