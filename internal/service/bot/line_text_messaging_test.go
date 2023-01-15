package bot

import (
	"encoding/json"
	"testing"

	"github.com/line/line-bot-sdk-go/linebot"
)

func TestGetTodo(t *testing.T) {
	lineId := LineID{
		UserID:  "test",
		GroupID: "test",
		RoomID:  "test",
	}
	getTodo(lineId)
}

func TestCreateTodo(t *testing.T) {
	templates := []interface{}{}
	templates = append(templates, linebot.NewTextMessage("test"))
	t.Logf("%v", templates)

	templatesJSONByte, err := json.Marshal(templates)
	if err != nil {
		t.Errorf("createTodo templatesJSONByte json.Marshal error: %v", err)
		return
	}
	t.Logf("%v", string(templatesJSONByte))
}
