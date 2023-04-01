package helper

import (
	"os"
	"path/filepath"
	"runtime"
)

func GetCurrentGoFilePath() string {
	goFilePath, _ := filepath.Abs(os.Args[0])
	return goFilePath
}

func GetFunctionName() string {
	// Skip 2 levels to get the caller function name
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	return funcName
}
