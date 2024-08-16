package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ernestngugi/todo/internal/db"
	"github.com/ernestngugi/todo/internal/entities"
	"github.com/ernestngugi/todo/internal/forms"
	"github.com/ernestngugi/todo/internal/repository"
	"github.com/ernestngugi/todo/internal/utils"
)

type (
	TodoController interface {
		CompleteTodo(ctx context.Context, dB db.DB, todoID int64) (*entities.Todo, error)
		CreateTodo(ctx context.Context, dB db.DB, form *forms.CreateTodoForm) (*entities.Todo, error)
		DeleteTodo(ctx context.Context, dB db.DB, todoID int64) error
		TodoByID(ctx context.Context, dB db.DB, todoID int64) (*entities.Todo, error)
		Todos(ctx context.Context, dB db.DB, filter *forms.Filter) (*entities.TodoList, error)
		UpdateTodo(ctx context.Context, dB db.DB, todoID int64, form *forms.UpdateTodoForm) (*entities.Todo, error)
	}

	todoController struct {
		todoRepository repository.TodoRepository
	}
)

func NewTestTodoController() *todoController {
	return &todoController{
		todoRepository: repository.NewTodoRepository(),
	}
}

func NewTodoController(todoRepository repository.TodoRepository) TodoController {
	return &todoController{
		todoRepository: todoRepository,
	}
}

func (s *todoController) TodoByID(ctx context.Context, dB db.DB, todoID int64) (*entities.Todo, error) {
	return s.todoRepository.TodoByID(ctx, dB, todoID)
}

func (s *todoController) CreateTodo(ctx context.Context, dB db.DB, form *forms.CreateTodoForm) (*entities.Todo, error) {

	err := utils.ValidateSingleName(form.Title)
	if err != nil {
		return &entities.Todo{}, err
	}

	todo := &entities.Todo{
		Title: form.Title,
	}

	if strings.TrimSpace(form.Description) != "" {
		todo.Description = form.Description
	}

	err = s.todoRepository.Save(ctx, dB, todo)
	if err != nil {
		return &entities.Todo{}, err
	}

	return todo, nil
}

func (s *todoController) UpdateTodo(ctx context.Context, dB db.DB, todoID int64, form *forms.UpdateTodoForm) (*entities.Todo, error) {

	todo, err := s.todoRepository.TodoByID(ctx, dB, todoID)
	if err != nil {
		return &entities.Todo{}, err
	}

	if form.Title != nil {
		err := utils.ValidateSingleName(*form.Title)
		if err != nil {
			return &entities.Todo{}, err
		}
		todo.Title = *form.Title
	}

	if form.Description != nil {
		if strings.TrimSpace(*form.Description) != "" {
			todo.Description = *form.Description
		}
	}

	err = s.todoRepository.Save(ctx, dB, todo)
	if err != nil {
		return &entities.Todo{}, err
	}

	return todo, nil
}

func (s *todoController) CompleteTodo(ctx context.Context, dB db.DB, todoID int64) (*entities.Todo, error) {

	todo, err := s.todoRepository.TodoByID(ctx, dB, todoID)
	if err != nil {
		return &entities.Todo{}, err
	}

	if todo.Completed {
		return &entities.Todo{}, fmt.Errorf("todo has been marked as complete")
	}

	timeNow := time.Now()

	todo.Completed = true
	todo.CompletedAt = &timeNow

	err = s.todoRepository.Save(ctx, dB, todo)
	if err != nil {
		return &entities.Todo{}, err
	}

	return todo, nil
}

func (s *todoController) DeleteTodo(ctx context.Context, dB db.DB, todoID int64) error {

	todo, err := s.todoRepository.TodoByID(ctx, dB, todoID)
	if err != nil {
		return err
	}

	if todo.Completed {
		return fmt.Errorf("cannot a todo that has been completed")
	}

	return s.todoRepository.DeleteTodo(ctx, dB, todo.ID)
}

func (s *todoController) Todos(ctx context.Context, dB db.DB, filter *forms.Filter) (*entities.TodoList, error) {

	todos, err := s.todoRepository.Todos(ctx, dB, filter)
	if err != nil {
		return &entities.TodoList{}, err
	}

	count, err := s.todoRepository.NumberOfTodos(ctx, dB, filter)
	if err != nil {
		return &entities.TodoList{}, err
	}

	todoList := &entities.TodoList{
		Todos:      todos,
		Pagination: entities.NewPagination(count, filter.Page, filter.Per),
	}

	return todoList, nil
}
