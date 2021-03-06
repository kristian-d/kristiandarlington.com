package ui

import (
	"github.com/kristian-d/kristiandarlington.com/internal/projectpath"
	"github.com/shurcooL/httpfs/union"
	"net/http"
	"path"
)

// Assets contains the project's assets.
var Assets = func() http.FileSystem {
	assetsPrefix := path.Join(projectpath.Root, "web/ui/")

	// if I ever need to filter out files from the fs, can look at "github.com/shurcooL/httpfs/filter"
	static := http.Dir(path.Join(assetsPrefix, "static"))
	templates := http.Dir(path.Join(assetsPrefix, "templates"))
	wellKnown := http.Dir(path.Join(assetsPrefix, ".well-known")) // this was used for SSL DCV

	return union.New(map[string]http.FileSystem{
		"/templates": templates,
		"/static":    static,
		"/.well-known": wellKnown,
	})
}()