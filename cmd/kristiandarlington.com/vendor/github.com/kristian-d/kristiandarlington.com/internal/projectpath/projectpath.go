package projectpath

import (
	"os"
	"path/filepath"
	"runtime"
)

var (
	// root dir of this project
	Root = getRoot()
)

func getRoot() string {
	env := os.Getenv("ENV")
	switch env {
	case "prod":
		return "/app/"
	case "local":
		var _, b, _, _ = runtime.Caller(0)
		return filepath.Join(filepath.Dir(b), "../..")
	default:
		var _, b, _, _ = runtime.Caller(0)
		return filepath.Join(filepath.Dir(b), "../..")
	}
}
