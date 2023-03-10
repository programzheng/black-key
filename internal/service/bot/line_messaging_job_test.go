package bot_test

import (
	"testing"

	serviceBot "github.com/programzheng/black-key/internal/service/bot"
)

func TestRefreshTodoByAfterPushDateTime(t *testing.T) {
	serviceBot.RefreshTodoByAfterPushDateTime()
}
