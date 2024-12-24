package handlers

import (
	"gochat/fiber/db"

	"github.com/gofiber/fiber/v2"
)

func FetchAllRooms(c *fiber.Ctx) error {
	database, databaseConnErr := db.ConnectToDB()
	if databaseConnErr != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": databaseConnErr.Error(),
		})
	}
	rooms, databaseGetAllRoomErr := db.GetAllRooms(database)
	if databaseGetAllRoomErr != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": databaseGetAllRoomErr.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"rooms": rooms,
	})
}
