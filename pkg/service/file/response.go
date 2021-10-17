package file

import (
	"black-key/pkg/filesystem"
)

func getResponseFilePath() string {
	return filesystem.Driver.GetHostURL()
}
