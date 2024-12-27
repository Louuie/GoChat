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
	html := `<table border="1">
		<tr>
			<th>Room UUID</th>
			<th>Room Number</th>
		</tr>`
	for _, room := range *rooms {
		html += `<tr>
				<td>` + room.UUID + `</td>
				<td><button class="room-btn" data-uuid="` + room.UUID + `" data-room-number="` + room.Room_Number + `">` + room.Room_Number + `</button></td>
			</tr>`
	}
	html += `</table>`
	return c.Type("html").SendString(html)
}
