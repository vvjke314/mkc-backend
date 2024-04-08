package app

import "github.com/gin-gonic/gin"

// CreateFile загружает файл на сервер
// @Summary Загрузить файл
// @Description Загружает файл на сервер
// @Tags file
// @Security 	 BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param project_id path string true "Идентификатор проекта"
// @Param file formData file true "Файл для загрузки"
// @Failure 500 {object} errorResponse
// @Failure 401 {obejct} errorResponse
// @Failure 403 {object} errorResponse
// @Router /project/{project_id} [post]
func (a *Application) CreateFile(c *gin.Context) {

}

// DeleteFile удаляет файл с сервера
// @Summary Удалить файл
// @Description Удаляет файл с сервера
// @Tags file
// @Security 	 BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param project_id path string true "Идентификатор проекта"
// @Param  data body ds.DeleteFileReq true "CHANGE IT"
// @Failure 500 {object} errorResponse
// @Failure 401 {obejct} errorResponse
// @Failure 403 {object} errorResponse
// @Router /project/{project_id} [delete]
func (a *Application) DeleteFile(c *gin.Context) {

}

// UpdateFileName меняет имя файла на сервере
// @Summary Изменяет имя файла
// @Description Изменяет имя файла на сервере
// @Tags file
// @Security 	 BearerAuth
// @Accept multipart/form-data
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
// @Failure 400 {object} ApiResponse{"error": "Bad request"}
// @Failure 500 {object} ApiResponse{"error": "Internal server error"}
// @Router /project/{project_id}/file/{file_name} [get]
func (a *Application) GetFile(c *gin.Context) {

}
