package config

import (
	"context"
	"http_server/internal/repository/postgres"
	task "http_server/internal/service"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckConnection(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := os.Getenv("DB_CONNECTION")
	return pgxpool.New(ctx, dsn)
}

func InitPostgresEnv(ctx context.Context, pool *pgxpool.Pool) (*task.TaskService, error) {

	if err := createTaskTable(ctx, pool); err != nil {
		return &task.TaskService{}, err
	}

	postgresTaskRepository := postgres.NewPostgresTaskRepository(pool)
	return task.NewTaskService(postgresTaskRepository), nil
}

func createTaskTable(ctx context.Context, pool *pgxpool.Pool) error {

	if _, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS tasks (
			id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title        TEXT        NOT NULL,
			description  TEXT        NOT NULL DEFAULT '',
			status       BOOLEAN     NOT NULL DEFAULT FALSE,
			create_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			completed_at TIMESTAMPTZ,
			version      INTEGER     NOT NULL DEFAULT 1
		);`); err != nil {
		return err
	}

	return nil
}
