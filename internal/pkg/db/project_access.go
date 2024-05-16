package db

import (
	"context"
	"fmt"
	"log"

	pgx "github.com/jackc/pgx/v5"
	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
)

// CreateParticipant добавляет участника в проект
func (r *Repo) CreateParticipant(p ds.ProjectAccess) error {
	query := "INSERT INTO project_access (id, project_id, customer_id, customer_access) VALUES ($1, $2, $3, $4)"
	_, err := r.pool.Exec(r.ctx, query, p.Id, p.ProjectId, p.CustomerId, p.CustomerAccess)
	if err != nil {
		return fmt.Errorf("[pgxpool.Pool.Exec] Can't exec query %w", err)
	}
	return nil
}

// UpdateParticipantAccess [unchecked]
// Изменяет статус участника в проекте {0 : Чтение, 1 : Запись + Чтение}
func (r *Repo) UpdateParticipantAccess(participantId string, mod int) error {
	query := "UPDATE project_access SET customer_access = $1 WHERE customer_id = $2"
	_, err := r.pool.Exec(r.ctx, query, mod, participantId)
	if err != nil {
		return fmt.Errorf("[*pgxpool.Pool.Exec] Can't exec query: %w", err)
	}

	return nil
}

// DeleteParticipant удаляет участника из проекта
func (r *Repo) DeleteParticipant(participantId, projectId string) error {
	query := "DELETE FROM project_access WHERE customer_id = $1 AND project_id = $2"
	_, err := r.pool.Exec(r.ctx, query, participantId, projectId)
	if err != nil {
		return fmt.Errorf("[pgxpool.Pool.Exec] Can't exec query %w", err)
	}

	return nil
}

// DeleteParticipants удаляет всех участников из проекта
func (r *Repo) DeleteParticipants(projectId string) error {
	query := "DELETE FROM project_access WHERE project_id = $1"
	_, err := r.pool.Exec(r.ctx, query, projectId)
	if err != nil {
		return fmt.Errorf("[pgxpool.Pool.Exec] Can't exec query %w", err)
	}
	return nil
}

// AccessControl проверяет имеет ли пользователь доступ к проекту
func (r *Repo) AccessControl(customerId, projectId string, access int) (bool, error) {
	// var pa ds.ProjectAccess

	// Проверка на владельца проекта
	var p ds.Project
	err := r.GetProjectById(projectId, &p)
	if err != nil {
		return false, fmt.Errorf("[db.GetProjectById] %w", err)
	}

	if p.OwnerId.String() == customerId {
		return true, nil
	}

	// Проверка на участника проекта
	query := "SELECT FROM project_access WHERE customer_id = $1 AND project_id = $2 AND customer_access >= $3"
	err = r.pool.QueryRow(r.ctx, query, customerId, projectId, access).Scan()
	if err != nil {
		if err == pgx.ErrNoRows {
			// Если запись отсутствует, возвращаем nil и nil ошибку
			return false, nil
		}
		// Если возникла другая ошибка, возвращаем nil и эту ошибку
		return false, fmt.Errorf("[pgxpool.Pool.QueryRow.Scan] Can't exec query %w", err)
	}

	// Возвращаем указатель на user и nil ошибку
	return true, nil
}

// GetParticipants возвращает всех участников проекта включая владельца проекта (первый в списке)
func (r *Repo) GetParticipants(projectId string) ([]ds.Customer, error) {
	var customers []ds.Customer

	var customer ds.Customer
	query := "SELECT owner_id FROM project WHERE id = $1"
	err := r.pool.QueryRow(r.ctx, query, projectId).Scan(&customer.Id)
	if err != nil {
		return nil, fmt.Errorf("[pgxpool.Pool.QueryRow] Can't exec query %w", err)
	}

	query = "SELECT * FROM customer WHERE id = $1"
	err = r.pool.QueryRow(r.ctx, query, customer.Id.String()).Scan(&customer.Id, &customer.FirstName, &customer.SecondName, &customer.Login, &customer.Password, &customer.Email, &customer.Type)
	if err != nil {
		return nil, fmt.Errorf("[pgxpool.Pool.QueryRow] Can't exec query %w", err)
	}
	customers = append(customers, customer)

	rows, err := r.pool.Query(context.Background(), `
	SELECT c.*
	FROM customer c
	JOIN project_access pa ON c.id = pa.customer_id
	WHERE pa.project_id = $1
	`, projectId)
	if err != nil {
		return nil, fmt.Errorf("[pgxpool.Pool.Query] Can't exec query %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var customer ds.Customer
		err := rows.Scan(&customer.Id, &customer.FirstName, &customer.SecondName, &customer.Login, &customer.Password, &customer.Email, &customer.Type)
		if err != nil {
			log.Fatalf("Scan error: %v\n", err)
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

// CheckParticipant проверяет добавлен ли уже участник в проект и если да, то возравщает ошибку
func (r *Repo) CheckParticipant(participantId, projectId string) error {
	pa := &ds.ProjectAccess{}
	query := "SELECT id FROM project_access WHERE project_id = $1 AND customer_id = $2"
	err := r.pool.QueryRow(r.ctx, query, projectId, participantId).Scan(&pa.Id)
	if err == pgx.ErrNoRows {
		return nil
	}
	return fmt.Errorf("such participant already in this project")
}
