package db

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
)

// CreateFile добавление информации о файле в БД
func (r *Repo) CreateFile(f ds.File) error {
	query := "INSERT INTO file (id, project_id, filename, extension, size, file_path, update_datetime) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.pool.Exec(r.ctx, query, f.Id, f.ProjectId, f.Filename, f.Extension, f.Size, f.FilePath, f.UpdateDatetime)
	if err != nil {
		return fmt.Errorf("[*pgxpool.Pool.Exec] Can't exec query: %w", err)
	}
	return nil
}

// DeleteFile удаление информации о файле из БД
func (r *Repo) DeleteFile(fileId string) error {
	query := "DELETE FROM file WHERE id = $1"
	_, err := r.pool.Exec(r.ctx, query, fileId)
	if err != nil {
		return fmt.Errorf("[pgxpool.Pool.Exec] Can't exec query: %w", err)
	}

	return nil
}

// UpdateFileName изменение названия файла и времени последнего изменения файла
func (r *Repo) UpdateFileName(fileId, fileName string) error {
	query := "UPDATE file SET filename = $1, update_datetime = $2 WHERE id = $3"
	_, err := r.pool.Exec(r.ctx, query, fileName, time.Now(), fileId)
	if err != nil {
		return fmt.Errorf("[*pgxpool.Pool.Exec] Can't exec query: %w", err)
	}

	return nil
}

// GetFileById получение ифнормации о файле через БД
func (r *Repo) GetFileById(fileId string, f *ds.File) error {
	query := "SELECT id, project_id, filename, extension, size, file_path, update_datetime FROM file WHERE id = $1"
	err := r.pool.QueryRow(r.ctx, query, fileId).Scan(&f.Id, &f.ProjectId, &f.Filename, &f.Extension, &f.Size, &f.FilePath, &f.UpdateDatetime)
	if err != nil {
		return fmt.Errorf("[*pgxpool.Pool.QueryRow] Can't exec query: %w", err)
	}

	return nil
}

// GetFileByName получение информации о файле по имени через БД
func (r *Repo) GetFileByName(filename, extension, projectId string, f *ds.File) error {
	query := "SELECT id, project_id, filename, extension, size, file_path, update_datetime FROM file WHERE filename = $1 AND extension = $2 AND project_id = $3"
	err := r.pool.QueryRow(r.ctx, query, filename, extension, projectId).Scan(&f.Id, &f.ProjectId, &f.Filename, &f.Extension, &f.Size, &f.FilePath, &f.UpdateDatetime)
	if err != nil {
		return fmt.Errorf("[*pgxpool.Pool.QueryRow] Can't exec query: %w", err)
	}

	return nil
}

// GetFile получение всех файлов проекта
func (r *Repo) GetFiles(projectId string) ([]ds.File, error) {
	var files []ds.File
	query := "SELECT id, project_id, filename, extension, size, file_path, update_datetime FROM file WHERE project_id = $1"

	rows, err := r.pool.Query(r.ctx, query, projectId)
	if err != nil {
		return files, fmt.Errorf("[*pgxpool.Pool.Query] Can't exec query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var f ds.File
		if err := rows.Scan(&f.Id, &f.ProjectId, &f.Filename, &f.Extension, &f.Size, &f.FilePath, &f.UpdateDatetime); err != nil {
			return files, fmt.Errorf("[pgx.Rows.Scan] Can't scan data: %w", err)
		}
		files = append(files, f)
	}
	if err := rows.Err(); err != nil {
		return files, fmt.Errorf("[pgx.Rows.Err] Error occured when iterating data: %w", err)
	}

	return files, nil
}

// DeleteFiles удаляет все файлы из проекта
func (r *Repo) DeleteFiles(projectId string) error {
	query := "DELETE FROM file WHERE project_id = $1"

	_, err := r.pool.Exec(r.ctx, query, projectId)
	if err != nil {
		return fmt.Errorf("[*pgxpool.Pool.QueryRow] Can't exec query %w", err)
	}

	return nil
}

// CheckFileExistence проверяет на существование имя файла
func (r *Repo) CheckFileExistence(filename, extension, projectId string) error {
	f := ds.File{}
	query := "SELECT id FROM file WHERE filename = $1 AND project_id = $2 AND extension = $3"
	err := r.pool.QueryRow(r.ctx, query, filename, projectId, extension).Scan(&f.Id)
	if err == pgx.ErrNoRows {
		return nil
	}
	return fmt.Errorf("such file already in this project")
}
