package db

import (
	"context"
	"fmt"

	//"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
	"github.com/vvjke314/mkc-backend/internal/pkg/dsn"
)

type Repo struct {
	ctx    context.Context
	config *pgxpool.Config
	pool   *pgxpool.Pool
}

func NewRepo() *Repo {
	return &Repo{}
}

// Init медот для инициализации конфига
func (r *Repo) Init() error {
	url, err := dsn.GetDSN()
	if err != nil {
		return fmt.Errorf("[dsn.GetDSN]: can't get data string name: %w", err)
	}

	pgxConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return fmt.Errorf("[pgxpool.ParseConfig]: can't parse config: %w", err)
	}

	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	r.ctx = context.Background()
	r.config = pgxConfig
	return nil
}

// Connect создание pool'a для подключения к БД
func (r *Repo) Connect() error {
	pgxConnPool, err := pgxpool.NewWithConfig(context.TODO(), r.config)
	if err != nil {
		return fmt.Errorf("[pgxpool.NewWithConfig]: Can't parse config: %w", err)
	}

	r.pool = pgxConnPool
	return nil
}

// Close закрытие pgxPool
func (r *Repo) Close() {
	r.pool.Close()
}
