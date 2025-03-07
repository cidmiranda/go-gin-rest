package main

import (
	"github.com/cidmiranda/go-fiber-ws/database"
	_ "github.com/cidmiranda/go-fiber-ws/docs"
	"github.com/cidmiranda/go-fiber-ws/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

// @title GO Fiber API
// @version 1.0
// @description API desenvolvida em Go, fiber, postgres e gorm
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath /
func main() {
	database.Connect()
	app := fiber.New()
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://example.com/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
	}))
	app.Use(logger.New())
	app.Use(cors.New())
	router.SetupRoutes(app)
	// handle unavailable route
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})
	app.Listen(":8080")
}
