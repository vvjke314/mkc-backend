package db

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
)

// SignUpCustomer
// Добавляет нового пользователя в БД
func (r *Repo) SignUpCustomer(c ds.Customer) error {
	query := "INSERT INTO customer (id, first_name, second_name, login, password, email, type) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.pool.Exec(r.ctx, query, c.Id, c.FirstName, c.SecondName, c.Login, c.Password, c.Email, c.Type)
	if err != nil {
		return fmt.Errorf("[pgxpool.Pool.Exec] Can't exec query %w", err)
	}
	return nil
}

// UpgradeCustomerStatus
// Повышение статуса клиента
func (r *Repo) UpgradeCustomerStatus(customerId string, status int) error {
	query := "UPDATE customer SET type = $1 WHERE id = $2"
	_, err := r.pool.Exec(r.ctx, query, status, customerId)
	if err != nil {
		return fmt.Errorf("[pgxpool.Pool.Exec] Can't exec query %w", err)
	}

	return nil
}

// GetCustomerByEmail
// Получаем id клиента через email
func (r *Repo) GetCustomerByEmail(customerEmail string, c *ds.Customer) error {
	query := "SELECT id, first_name, second_name, login, password, email, type FROM customer WHERE email = $1"
	err := r.pool.QueryRow(r.ctx, query, customerEmail).Scan(&c.Id, &c.FirstName, &c.SecondName, &c.Login, &c.Password, &c.Email, &c.Type)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Если запись отсутствует, возвращаем nil ошибку
			return nil
		}
		return fmt.Errorf("[pgxpool.Pool.QueryRow] Can't exec query %w", err)
	}

	return nil
}
