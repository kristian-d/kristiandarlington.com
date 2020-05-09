package ui

import (
	"github.com/shurcooL/httpfs/union"
	"net/http"
	"path"
)

// Assets contains the project's assets.
var Assets = func() http.FileSystem {
	assetsPrefix := "/app/web/ui/"

	// if I ever need to filter out files from the fs, can look at "github.com/shurcooL/httpfs/filter"
	static := http.Dir(path.Join(assetsPrefix, "static"))
	templates := http.Dir(path.Join(assetsPrefix, "templates"))

	return union.New(map[string]http.FileSystem{
		"/templates": templates,
		"/static":    static,
	})
}()