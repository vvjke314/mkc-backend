package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/vvjke314/mkc-backend/internal/pkg/dsn"
)

type Repo struct {
	db *sql.DB
}

func NewRepo() *Repo {
	return &Repo{}
}

// Init
// медот для инициализации структуры sql.DB
func (r *Repo) Init() error {
	url, err := dsn.GetDSN()
	if err != nil {
		return fmt.Errorf("[dsn.GetDSN]: Can't get data string name: %w", err)
	}

	db, err := sql.Open("pgx", url)
	if err != nil {
		return fmt.Errorf("[sql.Open]: Can't open database: %w", err)
	}

	r.db = db
	return nil
}

// Connect
// метод для подключения к БД
func (r *Repo) Connect(ctx context.Context) (*sql.Conn, error) {
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("[sql.Db.Conn]: Can't connect to database: %w", err)
	}

	return conn, err
}
