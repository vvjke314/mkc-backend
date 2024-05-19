package db

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
	"golang.org/x/crypto/bcrypt"
)

// SignUpCustomer добавляет нового пользователя в БД
func (r *Repo) SignUpCustomer(c ds.Customer) error {
	query := "INSERT INTO customer (id, first_name, second_name, login, password, email, type) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.pool.Exec(r.ctx, query, c.Id, c.FirstName, c.SecondName, c.Login, c.Password, c.Email, c.Type)
	if err != nil {
		return fmt.Errorf("[pgxpool.Pool.Exec] Can't exec query %w", err)
	}
	return nil
}

// UpgradeCustomerStatus повышение статуса клиента
func (r *Repo) UpgradeCustomerStatus(customerId string, status int) error {
	query := "UPDATE customer SET type = $1 WHERE id = $2"
	_, err := r.pool.Exec(r.ctx, query, status, customerId)
	if err != nil {
		return fmt.Errorf("[pgxpool.Pool.Exec] Can't exec query %w", err)
	}

	return nil
}

// GetCustomerByEmail получает id клиента через email
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

// GetCustomerByCredentials получает id клиента через credentials
func (r *Repo) GetCustomerByCredentials(custCredentials ds.LoginCustomerReq, c *ds.Customer) error {
	hashedPassword, err := r.GetCustomerPassword(custCredentials.Login)
	if err != nil {
		return fmt.Errorf("[GetCustomerPassword] can't exec query %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(custCredentials.Password)); err != nil {
		return fmt.Errorf("[CompareHashAndPassword] mismatch hash and password %w", err)
	}

	query := "SELECT id, first_name, second_name, login, password, email, type FROM customer WHERE login = $1 AND password = $2"
	err = r.pool.QueryRow(r.ctx, query, custCredentials.Login, hashedPassword).Scan(&c.Id, &c.FirstName, &c.SecondName, &c.Login, &c.Password, &c.Email, &c.Type)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Если запись отсутствует, возвращаем nil ошибку
			return err
		}
		return fmt.Errorf("[pgxpool.Pool.QueryRow] Can't exec query %w", err)
	}

	return nil
}

// GetCustomerPassword получает хэшированный пароль клиента по логину
func (r *Repo) GetCustomerPassword(custLogin string) (string, error) {
	var password string
	query := "SELECT password FROM customer WHERE login = $1"
	err := r.pool.QueryRow(r.ctx, query, custLogin).Scan(&password)

	if err != nil {
		if err == pgx.ErrNoRows {
			// Если запись отсутствует, возвращаем nil ошибку
			return "", err
		}
		return "", fmt.Errorf("[pgxpool.Pool.QueryRow] Can't exec query %w", err)
	}

	return password, nil
}
