package seed

import (
	"black-key/pkg/model"
)

func All() {
	CreateAdmin(model.DB, "admin", "admin")
}
