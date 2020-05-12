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

/*func httpsRedirectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		proto := req.Header.Get("x-forwarded-proto")
		if proto == "http" || proto == "HTTP" {
			http.Redirect(res, req, fmt.Sprintf("https://%s%s", req.Host, req.URL), http.StatusPermanentRedirect)
			return
		}

		next.ServeHTTP(res, req)
	})
}*/

func httpsRedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("https://%s%s", r.Host, r.URL), http.StatusPermanentRedirect)
}

func NewRouter(cfg *config.Config) *mux.Router {
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
			Handler(route.handler).
			Schemes("https")
	}

	// for serving static files
	router.PathPrefix("/static/").Handler(http.FileServer(ui.Assets)).Schemes("https")
	router.PathPrefix("/").Schemes("http").HandlerFunc(httpsRedirectHandler)
	return router
}
