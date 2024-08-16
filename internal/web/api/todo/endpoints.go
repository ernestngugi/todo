package todo

import (
	"github.com/ernestngugi/todo/internal/controller"
	"github.com/ernestngugi/todo/internal/db"
	"github.com/gin-gonic/gin"
)

func AddOpenEndpoints(
	r *gin.RouterGroup,
	dB db.DB,
	todoController controller.TodoController,
) {
	r.POST("/todo", createTodo(dB, todoController))
	r.GET("/todos", listTodo(dB, todoController))
	r.GET("/todo/:id", todoByID(dB, todoController))
	r.PUT("/todo/:id", updateTodo(dB, todoController))
	r.POST("/todo/:id", completeTodo(dB, todoController))
	r.DELETE("/todo/:id", deleteTodo(dB, todoController))
}
