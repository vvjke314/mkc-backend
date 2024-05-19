package db

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
)

// SignUpAdministrator добавляет администратора в БД
func (r *Repo) SignUpAdministrator(a ds.Administrator) error {
	query := "INSERT INTO administrator (id, name, email, password) VALUES ($1, $2, $3, $4)"
	_, err := r.pool.Exec(r.ctx, query, a.Id, a.Name, a.Email, a.Password)
	if err != nil {
		return fmt.Errorf("[pgxpool.Pool.Exec] Can't exec query: %w", err)
	}
	return nil
}

// SetAdministrator устанавливает к проекту администратора
func (r *Repo) SetAdministrator(adminId, projectId string) error {
	query := "UPDATE project SET admin_id = $1 WHERE id = $2"
	_, err := r.pool.Exec(r.ctx, query, adminId, projectId)
	if err != nil {
		return fmt.Errorf("[pgxpool.Pool.Exec] Can't exec query: %w", err)
	}

	return nil
}

// GetAdminId получаем id администратора через его имя и пароль
func (r *Repo) GetAdminId(adminName, adminPass string) (string, error) {
	var id string
	query := "SELECT id FROM administrator WHERE name = $1 AND password = $2"
	row := r.pool.QueryRow(r.ctx, query, adminName, adminPass)
	err := row.Scan(&id)
	if err != nil {
		return "", fmt.Errorf("[pgxpool.Pool.Exec] Can't exec query: %w", err)
	}

	return id, nil
}

// GetCustomerEmail получает почту создателя проекта
func (r *Repo) GetCustomerEmail(projectId string) (string, error) {
	var email string
	query := `
		SELECT c.email
		FROM Project p
		JOIN Customer c ON p.owner_id = c.id
		WHERE p.id = $1;
	`

	err := r.pool.QueryRow(r.ctx, query, projectId).Scan(&email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("no customer found for project ID: %s", projectId)
		}
		return "", fmt.Errorf("query row failed: %w", err)
	}

	return email, nil
}

// GetAllUnattachedProjects возращает массив проектов, у которых значение admin_id = 0
func (r *Repo) GetAllUnattachedProjects() ([]ds.Project, error) {
	var projects []ds.Project
	query := "SELECT id, owner_id, capacity, name, creation_date, admin_id FROM project WHERE admin_id IS NULL"
	rows, err := r.pool.Query(r.ctx, query)
	if err != nil {
		return nil, fmt.Errorf("[pgxpool.Pool.Exec] Can't exec query: %w", err)
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

// GetAllAttachedProjects возращает массив проектов, которые прикрепены к администратору
func (r *Repo) GetAllAttachedProjects(adminId string) ([]ds.Project, error) {
	var projects []ds.Project
	query := "SELECT id, owner_id, capacity, name, creation_date, admin_id FROM project WHERE admin_id = $1"
	rows, err := r.pool.Query(r.ctx, query, adminId)
	if err != nil {
		return nil, fmt.Errorf("[pgxpool.Pool.Exec] Can't exec query: %w", err)
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

// GetValidCredentials получаем корректные данные
func (r *Repo) GetValidCredentials(name string) (string, error) {
	var hashedPassword string
	err := r.pool.QueryRow(r.ctx, "SELECT password FROM administrator WHERE name=$1", name).Scan(&hashedPassword)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}
