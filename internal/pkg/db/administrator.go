package db

import (
	"fmt"

	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
)

// SignUpAdministrator
// Добавляет администратора в БД
func (r *Repo) SignUpAdministrator(a ds.Administrator) error {
	query := "INSERT INTO administrator (id, name, email, password) VALUES ($1, $2, $3, $4)"
	_, err := r.pool.Exec(r.ctx, query, a.Id, a.Name, a.Email, a.Password)
	if err != nil {
		return fmt.Errorf("[*pgxpool.Pool.Exec] Can't exec query: %w", err)
	}
	return nil
}

// SetAdministrator
// Устанавливает к проекту администратора
func (r *Repo) SetAdministrator(adminId, projectId string) error {
	query := "UPDATE project SET admin_id = $1 WHERE id = $2"
	_, err := r.pool.Exec(r.ctx, query, adminId, projectId)
	if err != nil {
		return fmt.Errorf("[*pgxpool.Pool.Exec] Can't exec query: %w", err)
	}

	return nil
}
