package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kristian-d/kristiandarlington.com/config"
	"github.com/kristian-d/kristiandarlington.com/web/ui"
	"net/http"
)

type route struct {
	method string
	endpoint string
	handler http.HandlerFunc
}

type routes []route

func httpsRedirectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Printf("HOST:\t" + req.URL.Host + "\t" + req.Host)

		proto := req.Header.Get("x-forwarded-proto")
		if proto == "http" || proto == "HTTP" {
			http.Redirect(res, req, fmt.Sprintf("https://%s%s", req.Host, req.URL), http.StatusPermanentRedirect)
			return
		}

		next.ServeHTTP(res, req)
	})
}

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
	return httpsRedirectMiddleware(router)
}
