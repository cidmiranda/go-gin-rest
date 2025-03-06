package router

import (
	"github.com/cidmiranda/go-fiber-ws/handler"
	"github.com/cidmiranda/go-fiber-ws/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {
	// grouping
	api := app.Group("/api")
	api.Post("/auth", handler.LoginUser)
	v1 := api.Group("/user")
	// routes
	v1.Get("/", middleware.JWTProtected, handler.GetAllUsers)
	v1.Get("/:id", middleware.JWTProtected, handler.GetSingleUser)
	v1.Post("/", handler.CreateUser)
	v1.Put("/:id", middleware.JWTProtected, handler.UpdateUser)
	v1.Delete("/:id", middleware.JWTProtected, handler.DeleteUserByID)
	v1.Get("/login/retricted", middleware.JWTProtected, handler.Restricted)
}
