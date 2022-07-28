package main

import (
	"log"

	"github.com/long2ice/fibers/security"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/long2ice/fibers"
)

func main() {
	app := fibers.New(NewSwagger(), fiber.Config{ErrorHandler: func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		if _, ok := err.(validator.ValidationErrors); ok {
			code = fiber.StatusBadRequest
		}
		err = c.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		return nil
	}})
	app.Use(
		logger.New(),
		recover.New(),
		cors.New(),
	)
	subApp := fibers.New(NewSwagger(), fiber.Config{})
	subApp.Get("/noModel", noModel)
	app.Mount("/sub", subApp)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "*",
		AllowHeaders:     "*",
		AllowCredentials: true,
	}))
	queryGroup := app.Group("/query", fibers.Tags("Query"))
	queryGroup.Get("/list", queryList)
	queryGroup.Get("/:id", queryPath)
	queryGroup.Delete("", query)

	app.Get("/noModel", noModel)

	bodyGroup := app.Group("/body", fibers.Tags("Body"), fibers.Security(&security.Bearer{}))
	bodyGroup.Post("/encoded", formEncode)
	bodyGroup.Post("/file", file)
	bodyGroup.Post("/json", body)

	log.Fatal(app.Listen(":8080"))
}
