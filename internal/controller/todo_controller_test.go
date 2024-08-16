package controller

import (
	"context"
	"testing"
	"time"

	"github.com/ernestngugi/todo/internal/db"
	"github.com/ernestngugi/todo/internal/forms"
	"github.com/ernestngugi/todo/internal/repository"
	"github.com/ernestngugi/todo/internal/testutils"
	. "github.com/smartystreets/goconvey/convey"
	"syreclabs.com/go/faker"
)

func TestTodoController(t *testing.T) {

	testDB := db.InitDB()
	defer testDB.Close()

	ctx := context.Background()

	Convey("TestTodoController", t, testutils.WithTestDB(ctx, testDB, func(ctx context.Context, dB db.DB) {

		todoController := NewTestTodoController()

		Convey("can create a todo", func() {

			form := &forms.CreateTodoForm{
				Title:       "test",
				Description: faker.Lorem().Paragraph(3),
			}

			todo, err := todoController.CreateTodo(ctx, dB, form)
			So(err, ShouldBeNil)

			So(todo.Title, ShouldEqual, form.Title)
			So(todo.Description, ShouldEqual, form.Description)
			So(todo.Completed, ShouldBeFalse)
			So(todo.CompletedAt, ShouldBeZeroValue)
			So(todo.CreatedAt, ShouldNotBeZeroValue)
			So(todo.UpdatedAt, ShouldNotBeZeroValue)
		})

		Convey("can complete a todo", func() {

			todo, err := repository.CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			completedTodo, err := todoController.CompleteTodo(ctx, dB, todo.ID)
			So(err, ShouldBeNil)

			So(completedTodo.Completed, ShouldBeTrue)
			So(completedTodo.CompletedAt, ShouldNotBeZeroValue)
			So(completedTodo.UpdatedAt, ShouldNotEqual, completedTodo.CreatedAt)
		})

		Convey("cannot mark an already completed todo as complete", func() {

			todo, err := repository.CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			timeNow := time.Now()
			todo.Completed = true
			todo.CompletedAt = &timeNow

			err = todoController.todoRepository.Save(ctx, dB, todo)
			So(err, ShouldBeNil)

			_, err = todoController.CompleteTodo(ctx, dB, todo.ID)
			So(err, ShouldNotBeNil)

			So(err.Error(), ShouldEqual, "todo has been marked as complete")
		})

		Convey("can update a todo", func() {

			todo, err := repository.CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			title := "newtest"

			form := &forms.UpdateTodoForm{
				Title: &title,
			}

			updatedTodo, err := todoController.UpdateTodo(ctx, dB, todo.ID, form)
			So(err, ShouldBeNil)

			So(updatedTodo.Title, ShouldEqual, "newtest")
			So(updatedTodo.UpdatedAt, ShouldNotEqual, updatedTodo.CreatedAt)
		})

		Convey("can delete a todo", func() {

			todo, err := repository.CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			err = todoController.DeleteTodo(ctx, dB, todo.ID)
			So(err, ShouldBeNil)

			_, err = todoController.todoRepository.TodoByID(ctx, dB, todo.ID)
			So(err, ShouldNotBeNil)

			So(err.Error(), ShouldContainSubstring, "sql: no rows in result set")
		})

		Convey("cannot delete a completed todo", func() {

			todo, err := repository.CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			timeNow := time.Now()
			todo.Completed = true
			todo.CompletedAt = &timeNow

			err = todoController.todoRepository.Save(ctx, dB, todo)
			So(err, ShouldBeNil)

			err = todoController.DeleteTodo(ctx, dB, todo.ID)
			So(err, ShouldNotBeNil)

			So(err.Error(), ShouldEqual, "cannot a todo that has been completed")
		})

		Convey("can list todos", func() {

			todo1, err := repository.CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			todo2, err := repository.CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			todos, err := todoController.Todos(ctx, dB, &forms.Filter{})
			So(err, ShouldBeNil)

			So(len(todos.Todos), ShouldEqual, 2)

			foundTodo1 := todos.Todos[0]
			foundTodo2 := todos.Todos[1]

			So(foundTodo1.ID, ShouldEqual, todo1.ID)
			So(foundTodo2.ID, ShouldEqual, todo2.ID)
			So(todos.Pagination.Count, ShouldEqual, 2)
		})

		Convey("can get todo by id", func() {

			todo, err := repository.CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			foundTodo, err := todoController.TodoByID(ctx, dB, todo.ID)
			So(err, ShouldBeNil)

			So(foundTodo.ID, ShouldEqual, todo.ID)
			So(foundTodo.Title, ShouldEqual, todo.Title)
			So(foundTodo.Description, ShouldEqual, todo.Description)
			So(foundTodo.Completed, ShouldBeFalse)
			So(foundTodo.CompletedAt, ShouldBeZeroValue)
			So(foundTodo.CreatedAt, ShouldNotBeZeroValue)
			So(foundTodo.UpdatedAt, ShouldNotBeZeroValue)
		})
	}))
}
