package web

import (
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

type routes []route

func NewRouter(cfg *config.Config) http.Handler {
	var myRoutes = routes{
		route{
			"GET",
			"/",
			index,
		},
		route{
			"GET",
			"/projects/",
			projects,
		},
		route{
			"GET",
			"/resume/",
			resume,
		},
		route{
			"GET",
			"/about/",
			about,
		},
		route{
			"GET",
			"/contact/",
			contact,
		},
		route{
			"POST",
			"/contact/",
			func (w http.ResponseWriter, r *http.Request) {
				contactSend(w, r, *cfg)
			},
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
