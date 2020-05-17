package docs

import (
	"github.com/kristian-d/kristiandarlington.com/internal/projectpath"
	"net/http"
	"path"
)

var Docs = http.Dir(path.Join(projectpath.Root, "docs"))
