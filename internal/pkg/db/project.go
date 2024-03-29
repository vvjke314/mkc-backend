package db

import (
	"fmt"

	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
)

// CreateProject
// Добавляет новый проект в БД
func (r *Repo) CreateProject(p ds.Project) error {
	query := "INSERT INTO project (id, owner_id, capacity, name, creation_date) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.pool.Exec(r.ctx, query, p.Id, p.OwnerId, p.Capacity, p.Name, p.CreationDate)
	if err != nil {
		return fmt.Errorf("[pgxpool.Pool.Exec] Can't exec query %w", err)
	}
	return nil
}

// DeleteProject
// Удаляет проект из БД
func (r *Repo) DeleteProject(projectId string) error {
	err := r.DeleteFiles(projectId)
	if err != nil {
		return fmt.Errorf("[db.DeleteFiles] %w", err)
	}

	err = r.DeleteNotes(projectId)
	if err != nil {
		return fmt.Errorf("[db.DeleteNotes] %w", err)
	}

	err = r.DeleteParticipants(projectId)
	if err != nil {
		return fmt.Errorf("[db.DeleteParticipants] %w", err)
	}

	query := "DELETE FROM project WHERE id = $1"
	_, err = r.pool.Exec(r.ctx, query, projectId)
	if err != nil {
		return fmt.Errorf("[pgxpool.Pool.Exec] Can't exec query %w", err)
	}

	return nil
}

// UpdateProjectName
// Изменение имени проекта в БД
func (r *Repo) UpdateProjectName(projectId, projectName string) error {
	query := "UPDATE project SET name = $1 WHERE id = $2"
	_, err := r.pool.Exec(r.ctx, query, projectName, projectId)
	if err != nil {
		return fmt.Errorf("[*pgxpool.Pool.Exec] Can't exec query: %w", err)
	}

	return nil
}

// GetProjectById
// Получает структуру проект по id проекта
func (r *Repo) GetProjectById(projectId string, p *ds.Project) error {
	query := "SELECT id, owner_id, capacity, name, creation_date FROM project WHERE id = $1"
	err := r.pool.QueryRow(r.ctx, query, projectId).Scan(&p.Id, &p.OwnerId, &p.Capacity, &p.Name, &p.CreationDate)
	if err != nil {
		return fmt.Errorf("[*pgxpool.Pool.QueryRow] Can't exec query: %w", err)
	}

	return nil
}

// GetProjects
// Возращает все проекты пользователя
func (r *Repo) GetProjects(customerId string) ([]ds.Project, error) {
	var projects []ds.Project
	query := "SELECT id, owner_id, capacity, name, creation_date, admin_id FROM project WHERE owner_id = $1"

	rows, err := r.pool.Query(r.ctx, query, customerId)
	if err != nil {
		return projects, fmt.Errorf("[*pgxpool.Pool.Query] Can't exec query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p ds.Project
		if err := rows.Scan(&p.Id, &p.OwnerId, &p.Capacity, &p.Name, &p.CreationDate, &p.AdminId); err != nil {
			return projects, fmt.Errorf("[pgx.Rows.Scan] Can't scan data: %w", err)
		}
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return projects, fmt.Errorf("[pgx.Rows.Err] Error occured when iterating data: %w", err)
	}

	return projects, nil
}
