package postgres

import (
	"context"
	"errors"
	"http_server/internal/domain"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTaskRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresTaskRepository(pool *pgxpool.Pool) *PostgresTaskRepository {
	return &PostgresTaskRepository{
		pool: pool,
	}
}

func (repo *PostgresTaskRepository) Create(ctx context.Context, task domain.Task) error {
	if _, err := repo.pool.Exec(ctx, `
		INSERT INTO tasks (id, title, description, status, create_at, completed_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.CreateAt,
		task.CompletedAt,
		task.Version,
	); err != nil {
		return err
	}

	return nil
}

func (repo *PostgresTaskRepository) GetById(ctx context.Context, id uuid.UUID) (domain.Task, error) {

	var t domain.Task

	if err := repo.pool.QueryRow(ctx, `SELECT id, title, description, status, create_at, completed_at, version FROM tasks WHERE id = $1;`, id).Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.CreateAt,
		&t.CompletedAt,
		&t.Version); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, domain.ErrTaskNotFound
		}
		return domain.Task{}, err
	}

	return t, nil
}

func (repo *PostgresTaskRepository) List(ctx context.Context) ([]domain.Task, error) {

	rows, err := repo.pool.Query(ctx, `SELECT id, title, description, status, create_at, completed_at, version FROM tasks;`)
	if err != nil {
		return make([]domain.Task, 0), err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Task])
}

func (repo *PostgresTaskRepository) ListByCompletedStatus(ctx context.Context, status bool) ([]domain.Task, error) {

	rows, err := repo.pool.Query(ctx, `SELECT id, title, description, status, create_at, completed_at, version FROM tasks WHERE status = $1;`, status)
	if err != nil {
		return make([]domain.Task, 0), err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Task])
}

func (repo *PostgresTaskRepository) Update(ctx context.Context, task domain.Task) error {

	oldVersion := task.Version - 1

	result, err := repo.pool.Exec(ctx, `
		UPDATE tasks
		SET title = $2,
			description = $3,
			status = $4,
			completed_at = $5,
			version = $6
		WHERE id = $1 AND version = $7;`,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.CompletedAt,
		task.Version,
		oldVersion,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrTaskConflict
	}

	return nil
}

func (repo *PostgresTaskRepository) Delete(ctx context.Context, id uuid.UUID) error {

	result, err := repo.pool.Exec(ctx, `
		DELETE FROM tasks
		WHERE id = $1;
		`, id,
	)
	if err != nil {
		if result.RowsAffected() == 0 {
			return domain.ErrTaskNotFound
		}
		return err
	}

	return nil
}

func (repo *PostgresTaskRepository) GetVersionById(ctx context.Context, id uuid.UUID) (int, error) {

	var version int

	if err := repo.pool.QueryRow(ctx,
		`SELECT version FROM tasks WHERE id = $1;`, id).Scan(&version); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, domain.ErrTaskNotFound
		}
		return 0, err
	}

	return version, nil
}
