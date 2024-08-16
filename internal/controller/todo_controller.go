package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ernestngugi/todo/internal/db"
	"github.com/ernestngugi/todo/internal/entities"
	"github.com/ernestngugi/todo/internal/forms"
	"github.com/ernestngugi/todo/internal/providers"
	"github.com/ernestngugi/todo/internal/repository"
	"github.com/ernestngugi/todo/internal/utils"
)

const (
	todoKeyPrefix = "todo:todo-key:%v"
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
		cacheController CacheController
		todoRepository  repository.TodoRepository
	}
)

func NewTestTodoController(
	redisProvider providers.Redis,
) *todoController {
	cacheController := NewTestCacheController(redisProvider)
	return &todoController{
		cacheController: cacheController,
		todoRepository:  repository.NewTodoRepository(),
	}
}

func NewTodoController(
	cacheController CacheController,
	todoRepository repository.TodoRepository,
) TodoController {
	return &todoController{
		cacheController: cacheController,
		todoRepository:  todoRepository,
	}
}

func (s *todoController) TodoByID(ctx context.Context, dB db.DB, todoID int64) (*entities.Todo, error) {
	exist, err := s.cacheController.Exists(s.generateCacheKey(todoID))
	if err != nil {
		return &entities.Todo{}, err
	}

	var todo *entities.Todo

	if exist {

		err = s.cacheController.GetCachedValue(s.generateCacheKey(todoID), &todo)
		if err != nil {
			return &entities.Todo{}, err
		}

		return todo, nil
	}

	todo, err = s.todoRepository.TodoByID(ctx, dB, todoID)
	if err != nil {
		return &entities.Todo{}, err
	}

	return todo, nil
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

	err = s.cacheTodo(todo)
	if err != nil {
		return &entities.Todo{}, err
	}

	return todo, nil
}

func (s *todoController) UpdateTodo(ctx context.Context, dB db.DB, todoID int64, form *forms.UpdateTodoForm) (*entities.Todo, error) {

	exist, err := s.cacheController.Exists(s.generateCacheKey(todoID))
	if err != nil {
		return &entities.Todo{}, err
	}

	var todo *entities.Todo

	if exist {

		err = s.cacheController.GetCachedValue(s.generateCacheKey(todoID), &todo)
		if err != nil {
			return &entities.Todo{}, err
		}
	} else {

		todo, err = s.todoRepository.TodoByID(ctx, dB, todoID)
		if err != nil {
			return &entities.Todo{}, err
		}

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

	err = s.removeFromCache(todo.ID)
	if err != nil {
		return &entities.Todo{}, err
	}

	err = s.cacheTodo(todo)
	if err != nil {
		return &entities.Todo{}, err
	}

	return todo, nil
}

func (s *todoController) CompleteTodo(ctx context.Context, dB db.DB, todoID int64) (*entities.Todo, error) {

	exist, err := s.cacheController.Exists(s.generateCacheKey(todoID))
	if err != nil {
		return &entities.Todo{}, err
	}

	var todo *entities.Todo

	if exist {

		err = s.cacheController.GetCachedValue(s.generateCacheKey(todoID), &todo)
		if err != nil {
			return &entities.Todo{}, err
		}
	} else {

		todo, err = s.todoRepository.TodoByID(ctx, dB, todoID)
		if err != nil {
			return &entities.Todo{}, err
		}

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

	err = s.removeFromCache(todo.ID)
	if err != nil {
		return &entities.Todo{}, err
	}

	err = s.cacheTodo(todo)
	if err != nil {
		return &entities.Todo{}, err
	}

	return todo, nil
}

func (s *todoController) DeleteTodo(ctx context.Context, dB db.DB, todoID int64) error {

	exist, err := s.cacheController.Exists(s.generateCacheKey(todoID))
	if err != nil {
		return err
	}

	var todo *entities.Todo

	if exist {

		err = s.cacheController.GetCachedValue(s.generateCacheKey(todoID), &todo)
		if err != nil {
			return err
		}
	} else {

		todo, err = s.todoRepository.TodoByID(ctx, dB, todoID)
		if err != nil {
			return err
		}

	}

	if todo.Completed {
		return fmt.Errorf("cannot a todo that has been completed")
	}

	err = s.removeFromCache(todo.ID)
	if err != nil {
		return err
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

func (s *todoController) generateCacheKey(todoID int64) string {
	return fmt.Sprintf(todoKeyPrefix, todoID)
}

func (s *todoController) cacheTodo(todo *entities.Todo) error {
	return s.cacheController.CacheValue(s.generateCacheKey(todo.ID), todo)
}

func (s *todoController) removeFromCache(todoID int64) error {
	return s.cacheController.RemoveFromCache(s.generateCacheKey(todoID))
}
