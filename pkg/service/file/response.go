package file

import (
	"github.com/programzheng/black-key/pkg/filesystem"
)

func getResponseFilePath() string {
	return filesystem.Driver.GetHostURL()
}
