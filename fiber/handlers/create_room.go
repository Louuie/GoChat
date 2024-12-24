package handlers

import (
	"gochat/fiber/db"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/rand"
)

func CreateRoom(c *fiber.Ctx) error {
	// Generate random UUIDs
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	generated_uuid := RandStringBytes(6, letterBytes)

	database, dbConnErr := db.ConnectToDB()
	if dbConnErr != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": dbConnErr.Error(),
		})
	}
	dbCreateTableErr := db.CreateTable(database)
	if dbCreateTableErr != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": dbCreateTableErr.Error(),
		})
	}
	UUID, room_number, roomCreateErr := db.CreateRoom(database, generated_uuid)
	if roomCreateErr != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": roomCreateErr.Error(),
		})
	}
	database.Close()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"uuid":        UUID,
		"room_number": room_number,
		"status":      "Good!",
	})
}
func RandStringBytes(n int, letterBytes string) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
