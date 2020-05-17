package web

import (
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"github.com/kristian-d/kristiandarlington.com/config"
	"github.com/kristian-d/kristiandarlington.com/web/ui"
	"net/http"
	"os"
)

type route struct {
	method string
	endpoint string
	handler http.HandlerFunc
}

type handler struct {
	cfg *config.Config
	logger *logger.Logger
}

type routes []route

func handlerize(fn func (http.ResponseWriter, *http.Request, interface{})) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		fn(w, r, nil)
	}
}

func NewRouter(cfg *config.Config, logger *logger.Logger) http.Handler {
	var h = &handler{
		cfg: cfg,
		logger: logger,
	}

	var myRoutes = routes{
		route{
			"GET",
			"/",
			handlerize(h.index),
		},
		route{
			"GET",
			"/projects/",
			handlerize(h.projects),
		},
		route{
			"GET",
			//"/projects/{name:^[a-z_]+\\.pdf$}", TODO: investigate
			"/projects/{filename}",
			handlerize(h.projectReports),
		},
		route{
			"GET",
			"/resume/",
			handlerize(h.resume),
		},
		route{
			"GET",
			"/about/",
			handlerize(h.about),
		},
		route{
			"GET",
			"/contact/",
			handlerize(h.contact),
		},
		route{
			"POST",
			"/contact/",
			handlerize(h.contactSend),
		},
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range myRoutes {
		router.
			Methods(route.method).
			Path(route.endpoint).
			Handler(route.handler)
	}

	// for serving static files
	router.PathPrefix("/static/").Handler(http.FileServer(ui.Assets))

	switch os.Getenv("ENV") {
	case "prod":
		return httpsRedirectMiddleware(router)
	case "local":
		return router
	default:
		return router
	}
}
