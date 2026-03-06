package repository

import (
	"Ozoi/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTask(pool *pgxpool.Pool, title string, completed bool, description string, userID string) (*models.OzoiTask, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO ozoi_tasks (title, completed, description, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, completed, description, created_at, updated_at, user_id
`

	var repoTask models.OzoiTask

	var err error = pool.QueryRow(ctx, query, title, completed, description, userID).Scan(
		&repoTask.ID,
		&repoTask.Title,
		&repoTask.Completed,
		&repoTask.Description,
		&repoTask.CreatedAt,
		&repoTask.UpdatedAt,
		&repoTask.UserID,
	)
	if err != nil {
		return nil, err
	}
	return &repoTask, nil
}

func GetAllTasks(pool *pgxpool.Pool, userID string) ([]*models.OzoiTask, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query = `
	SELECT * FROM ozoi_tasks
	WHERE user_id = $1
	ORDER BY created_at DESC
	`

	var rows, err = pool.Query(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var repoTask []*models.OzoiTask

	for rows.Next() {
		var item models.OzoiTask

		err = rows.Scan(
			&item.ID,
			&item.Title,
			&item.Completed,
			&item.Description,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.UserID,
		)

		if err != nil {
			return nil, err
		}

		repoTask = append(repoTask, &item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return repoTask, nil
}

func GetTaskByID(pool *pgxpool.Pool, id int, userID string) (*models.OzoiTask, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query = `
	SELECT * FROM ozoi_tasks
	WHERE id = $1 AND user_id = $2
	`

	var repoTask models.OzoiTask

	var err error = pool.QueryRow(ctx, query, id, userID).Scan(
		&repoTask.ID,
		&repoTask.Title,
		&repoTask.Completed,
		&repoTask.Description,
		&repoTask.CreatedAt,
		&repoTask.UpdatedAt,
		&repoTask.UserID,
	)

	if err != nil {
		return nil, err
	}
	return &repoTask, nil
}

func UpdateTaskByID(pool *pgxpool.Pool, id int, title string, description string, completed *bool, userID string) (*models.OzoiTask, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query = `
	UPDATE ozoi_tasks
	SET
		title       = COALESCE(NULLIF($1, ''), title),
		completed   = COALESCE($2, completed),
		description = COALESCE(NULLIF($3, ''), description),
		updated_at  = CURRENT_TIMESTAMP
	WHERE id = $4 AND user_id = $5
	RETURNING id, title, completed, description, created_at, updated_at, user_id
`

	var updated models.OzoiTask

	err := pool.QueryRow(ctx, query, title, completed, description, id, userID).Scan(
		&updated.ID,
		&updated.Title,
		&updated.Completed,
		&updated.Description,
		&updated.CreatedAt,
		&updated.UpdatedAt,
		&updated.UserID,
	)

	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func DeleteTaskByID(pool *pgxpool.Pool, id int, userID string) error {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query = `
	DELETE from ozoi_tasks
	WHERE id = $1 AND user_id = $2
`
	commandTag, err := pool.Exec(ctx, query, id, userID)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("Task %d is not found", id)
	}

	return nil
}
