package router

import (
	"net/http"

	"github.com/ernestngugi/todo/internal/controller"
	"github.com/ernestngugi/todo/internal/db"
	"github.com/ernestngugi/todo/internal/repository"
	"github.com/ernestngugi/todo/internal/web/api/todo"
	"github.com/gin-gonic/gin"
)

type AppRouter struct {
	*gin.Engine
}

func BuildRouter(
	dB db.DB,
) *AppRouter {

	router := gin.Default()

	appRouter := router.Group("/v1")

	todoRepository := repository.NewTodoRepository()

	todoController := controller.NewTodoController(todoRepository)

	todo.AddOpenEndpoints(appRouter, dB, todoController)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error_message": "Endpoint not found"})
	})

	return &AppRouter{
		router,
	}
}
