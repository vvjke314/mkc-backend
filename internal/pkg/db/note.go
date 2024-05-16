package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/vvjke314/mkc-backend/internal/pkg/ds"
)

// CreateNote
// Добавление информации о заметке в БД
func (r *Repo) CreateNote(n ds.Note) error {
	query := "INSERT INTO note (id, project_id, title, content, update_datetime, deadline) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := r.pool.Exec(r.ctx, query, n.Id, n.ProjectId, n.Title, n.Content, n.UpdateDatetime, n.Deadline)
	if err != nil {
		return fmt.Errorf("[pgxpool.Pool.Exec] can't exec query %w", err)
	}
	return nil
}

// DeleteNote
// Удаление информации о заметке из БД
func (r *Repo) DeleteNote(noteId string) error {
	query := "DELETE FROM note WHERE id = $1"
	_, err := r.pool.Exec(r.ctx, query, noteId)
	if err != nil {
		return fmt.Errorf("[pgxpool.Pool.Exec] can't exec query %w", err)
	}

	return nil
}

// UpdateNoteName
// Изменение названия заметки
func (r *Repo) UpdateNoteName(noteId, noteName string) error {
	query := "UPDATE note SET title = $1, update_datetime = $2 WHERE id = $3"
	_, err := r.pool.Exec(r.ctx, query, noteName, time.Now(), noteId)
	if err != nil {
		return fmt.Errorf("[pgxpool.Pool.Exec] can't exec query %w", err)
	}

	return nil
}

// UpdateNoteDeadLine
// Изменение дедлайна заметки
func (r *Repo) UpdateNoteDeadLine(noteId string, deadline time.Time) error {
	query := "UPDATE note SET deadline = $1, update_datetime = $2 WHERE id = $3"
	_, err := r.pool.Exec(r.ctx, query, deadline, time.Now(), noteId)
	if err != nil {
		return fmt.Errorf("[pgxpool.Pool.Exec] can't exec query %w", err)
	}

	return nil
}

// GetNoteById
// Получение ифнормации о заметке через БД
func (r *Repo) GetNoteById(noteId string, n *ds.Note) error {
	query := "SELECT id, project_id, title, content, update_datetime, deadline FROM note WHERE id = $1"
	err := r.pool.QueryRow(r.ctx, query, noteId).Scan(&n.Id, &n.ProjectId, &n.Title, &n.Content, &n.UpdateDatetime, &n.Deadline)
	if err != nil {
		if err == pgx.ErrNoRows {
			return err
		}
		return fmt.Errorf("[pgxpool.Pool.Exec] can't exec query %w", err)
	}

	return nil
}

// GetNoteByName
// Получение заметки по его имени в конкретном проекте
func (r *Repo) GetNoteByName(noteTitle, projectId string) error {
	var n string
	query := "SELECT id FROM note WHERE project_id = $1 AND title = $2"
	err := r.pool.QueryRow(r.ctx, query, projectId, noteTitle).Scan(&n)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("no such notes")
		}
		return fmt.Errorf("[pgxpool.Pool.Exec] can't exec query %w", err)
	}

	return nil
}

// GetNotes
// Получение всех файлов проекта
func (r *Repo) GetNotes(projectId string) ([]ds.Note, error) {
	var notes []ds.Note
	query := "SELECT id, project_id, title, content, update_datetime, deadline FROM note WHERE project_id = $1"

	rows, err := r.pool.Query(r.ctx, query, projectId)
	if err != nil {
		return notes, fmt.Errorf("[pgxpool.Pool.Query] can't exec query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var n ds.Note
		if err := rows.Scan(&n.Id, &n.ProjectId, &n.Title, &n.Content, &n.UpdateDatetime, &n.Deadline); err != nil {
			return notes, fmt.Errorf("[pgx.Rows.Scan] can't scan data: %w", err)
		}
		notes = append(notes, n)
	}
	if err := rows.Err(); err != nil {
		return notes, fmt.Errorf("[pgx.Rows.Err] error occured when iterating data: %w", err)
	}

	return notes, nil
}

// DeleteNotes
// Удаляет все уведомления в проекте
func (r *Repo) DeleteNotes(projectId string) error {
	query := "DELETE FROM note WHERE project_id = $1"

	_, err := r.pool.Exec(r.ctx, query, projectId)
	if err != nil {
		return fmt.Errorf("[pgxpool.Pool.QueryRow] can't exec query %w", err)
	}

	return nil
}
