package filesystem

import (
	"black-key/config"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type FileSystem interface {
	Check()
	GetSystem() string
	GetPath() string
	Upload(*gin.Context, *multipart.FileHeader) error
	GetHostURL() string
}

var Driver FileSystem

func init() {
	system := config.Cfg.GetString("FILESYSTEM_DRIVER")
	switch system {
	case "local":
		Driver = Local{
			System: config.Cfg.GetString("FILESYSTEM_DRIVER"),
			Path:   config.Cfg.GetString("FILESYSTEM_LOCAL_PATH"),
		}
	}
}
