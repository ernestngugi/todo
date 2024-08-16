package repository

import (
	"context"

	"github.com/ernestngugi/todo/internal/db"
	"github.com/ernestngugi/todo/internal/entities"
)

func CreateTodo(ctx context.Context, dB db.DB) (*entities.Todo, error) {
	todo := entities.BuildTodo()
	err := NewTodoRepository().Save(ctx, dB, todo)
	return todo, err
}
