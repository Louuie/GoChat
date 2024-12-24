package fiber

import (
	"gochat/fiber/handlers"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func StartServer() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Access-Control-Allow-Origin, Content-Type, Origin, Accept",
	}), logger.New())

	app.Post("/create-room", handlers.CreateRoom)
	app.Get("/fetch-rooms", handlers.FetchAllRooms)
	log.Fatal(app.Listen(":" + "5050"))
}
