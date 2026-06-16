package tasks_transport_http

import (
	"time"

	"github.com/Tim73916/go-todoapp/internal/core/domain"
)

type TaskDTOResponse struct {
	ID           int        `json:"id" example:"15"`
	Version      int        `json:"version" example:"3"`
	Title        string     `json:"title" example:"Homework"`
	Description  *string    `json:"description" example:"Make mathematics homework until Friday"`
	Completed    bool       `json:"completed" example:"false"`
	CreatedAt    time.Time  `json:"created_at" example:"2026-01-24T12:28:00Z"`
	CompletedAt  *time.Time `json:"completed_at" example:"null"`
	AuthorUserID int        `json:"author_user_id" example:"5"`
}

func taskDTOFromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func taskDTOsFromDomains(tasks []domain.Task) []TaskDTOResponse {
	dtos := make([]TaskDTOResponse, len(tasks))

	for i, task := range tasks {
		dtos[i] = taskDTOFromDomain(task)
	}

	return dtos
}
