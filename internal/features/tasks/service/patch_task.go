package tasks_service

import (
	"context"
	"fmt"
	"time"

	"github.com/Tim73916/go-todoapp/internal/core/domain"
)

func (s *TasksService) PatchTask(
	ctx context.Context,
	id int,
	patch domain.TaskPatch,
) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}

	if err := task.ApplyPatch(patch); err != nil {
		return domain.Task{}, fmt.Errorf("apply task patch: %w", err)
	}

	if patch.Completed.Set {
		if task.Completed {
			now := time.Now()
			task.CompletedAt = &now
		} else {
			task.CompletedAt = nil
		}
	}

	patchedTask, err := s.tasksRepository.PatchTask(ctx, id, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("patch task: %w", err)
	}

	return patchedTask, nil
}
