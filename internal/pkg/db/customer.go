package db

import (
	"context"
	"fmt"

	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
)

// SignUpCustomer
// Добавляет нового пользователя в БД
func (r *Repo) SignUpCustomer(c ds.Customer) error {
	query := "INSERT INTO Customer (id, first_name, second_name, login, password, email, type) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.pool.Exec(context.Background(), query, c.Id, c.FirstName, c.SecondName, c.Login, c.Password, c.Email, c.Type)
	if err != nil {
		return fmt.Errorf("[*pgxpool.Pool.Exec] Can't exec query: %w", err)
	}
	return nil
}
