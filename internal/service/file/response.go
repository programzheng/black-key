package file

import (
	"github.com/programzheng/black-key/internal/filesystem"
)

func getResponseFilePath() string {
	return filesystem.Driver.GetHostURL()
}
