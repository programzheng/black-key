package bot

import "testing"

func TestGetTodo(t *testing.T) {
	lineId := LineID{
		UserID:  "test",
		GroupID: "test",
		RoomID:  "test",
	}
	getTodo(lineId)
}
