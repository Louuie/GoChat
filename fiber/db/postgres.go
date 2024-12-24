package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Room struct {
	UUID        string `json:"uuid,omitempty"`
	Room_Number string `json:"room_number,omitempty"`
}

func ConnectToDB() (*sql.DB, error) {
	fmt.Println(os.Getenv("GO_CHAT_DB"))
	db, err := sql.Open("postgres", os.Getenv("GO_CHAT_DB"))
	if err != nil {
		return nil, err
	}
	ping := db.Ping()
	if ping != nil {
		return nil, ping
	}
	return db, nil
}
func CreateTable(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS rooms (UUID VARCHAR NOT NULL, room_number INTEGER NOT NULL, PRIMARY KEY(UUID))")
	if err, ok := err.(*pq.Error); ok {
		return errors.New(err.Code.Name())
	}
	return nil
}
func CreateRoom(db *sql.DB, UUID string) (string, int, error) {
	var room_nuber = 1
	_, err := db.Exec("INSERT into rooms VALUES ($1, $2)", UUID, room_nuber)
	if err, ok := err.(*pq.Error); ok {
		// 23505: unique_violation
		if err.Code.Name() == "unique_violation" {
			return "", 0, errors.New("that song is already in the queue")
		}
	}
	return UUID, room_nuber, nil
}

func GetAllRooms(db *sql.DB) (*[]Room, error) {
	roomsQuery, err := db.Query("SELECT * FROM rooms")
	if err, ok := err.(*pq.Error); ok {
		return nil, err
	}
	rooms := make([]Room, 0)
	for roomsQuery.Next() {
		room := Room{}
		err := roomsQuery.Scan(&room.UUID, &room.Room_Number)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return &rooms, nil
}
