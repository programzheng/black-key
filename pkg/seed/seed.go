package seed

import (
	"github.com/programzheng/black-key/pkg/model"
)

func All() {
	CreateAdmin(model.DB, "admin", "admin")
}
