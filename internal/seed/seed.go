package seed

import (
	"github.com/programzheng/black-key/internal/model"
)

func All() {
	CreateAdmin(model.DB, "admin", "admin")
}
