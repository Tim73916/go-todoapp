package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Tim73916/go-todoapp/internal/core/domain"
)

func (r *TasksRepository) GetTasks(
	ctx context.Context,
	userID *int,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	defaultLimit := 100
	defaultOffset := 0

	limitVal := defaultLimit
	if limit != nil {
		limitVal = *limit
	}

	offsetVal := defaultOffset
	if offset != nil {
		offsetVal = *offset
	}

	query := `
		SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
		FROM todoapp.tasks
		WHERE ($1::int IS NULL OR author_user_id = $1)
		ORDER BY id ASC
		LIMIT $2
		OFFSET $3;
	`

	args := []any{userID, limitVal, offsetVal}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel

	for rows.Next() {
		var taskModel TaskModel

		err := rows.Scan(
			&taskModel.ID,
			&taskModel.Version,
			&taskModel.Title,
			&taskModel.Description,
			&taskModel.Completed,
			&taskModel.CreatedAt,
			&taskModel.CompletedAt,
			&taskModel.AuthorUserID,
		)
		if err != nil {
			return nil, fmt.Errorf("scan tasks: %w", err)
		}

		taskModels = append(taskModels, taskModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return taskDomainFromModels(taskModels), nil
}
