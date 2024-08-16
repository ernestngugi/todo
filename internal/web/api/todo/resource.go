package todo

import (
	"net/http"
	"strconv"

	"github.com/ernestngugi/todo/internal/apperror"
	"github.com/ernestngugi/todo/internal/controller"
	"github.com/ernestngugi/todo/internal/db"
	"github.com/ernestngugi/todo/internal/forms"
	"github.com/ernestngugi/todo/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

func createTodo(
	dB db.DB,
	todoController controller.TodoController,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		var form forms.CreateTodoForm

		err := c.BindJSON(&form)
		if err != nil {
			appError := apperror.Wrap(err)
			webutils.HandleError(c, appError)
			return
		}

		todo, err := todoController.CreateTodo(c.Request.Context(), dB, &form)
		if err != nil {
			appError := apperror.Wrap(err)
			webutils.HandleError(c, appError)
			return
		}

		c.JSON(http.StatusOK, todo)
	}
}

func completeTodo(
	dB db.DB,
	todoController controller.TodoController,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			appError := apperror.Wrap(err)
			webutils.HandleError(c, appError)
			return
		}

		todo, err := todoController.CompleteTodo(c.Request.Context(), dB, todoID)
		if err != nil {
			appError := apperror.Wrap(err)
			webutils.HandleError(c, appError)
			return
		}

		c.JSON(http.StatusOK, todo)
	}
}

func updateTodo(
	dB db.DB,
	todoController controller.TodoController,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		var form forms.UpdateTodoForm

		err := c.BindJSON(&form)
		if err != nil {
			appError := apperror.Wrap(err)
			webutils.HandleError(c, appError)
			return
		}

		todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			appError := apperror.Wrap(err)
			webutils.HandleError(c, appError)
			return
		}

		todo, err := todoController.UpdateTodo(c.Request.Context(), dB, todoID, &form)
		if err != nil {
			appError := apperror.Wrap(err)
			webutils.HandleError(c, appError)
			return
		}

		c.JSON(http.StatusOK, todo)
	}
}

func todoByID(
	dB db.DB,
	todoController controller.TodoController,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			appError := apperror.Wrap(err)
			webutils.HandleError(c, appError)
			return
		}

		todo, err := todoController.TodoByID(c.Request.Context(), dB, todoID)
		if err != nil {
			appError := apperror.Wrap(err)
			webutils.HandleError(c, appError)
			return
		}

		c.JSON(http.StatusOK, todo)
	}
}

func deleteTodo(
	dB db.DB,
	todoController controller.TodoController,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			appError := apperror.Wrap(err)
			webutils.HandleError(c, appError)
			return
		}

		err = todoController.DeleteTodo(c.Request.Context(), dB, todoID)
		if err != nil {
			appError := apperror.Wrap(err)
			webutils.HandleError(c, appError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func listTodo(
	dB db.DB,
	todoController controller.TodoController,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		filter, err := webutils.FilterFromContext(c)
		if err != nil {
			appError := apperror.Wrap(err)
			webutils.HandleError(c, appError)
			return
		}

		todos, err := todoController.Todos(c.Request.Context(), dB, filter)
		if err != nil {
			appError := apperror.Wrap(err)
			webutils.HandleError(c, appError)
			return
		}

		c.JSON(http.StatusOK, todos)
	}
}
