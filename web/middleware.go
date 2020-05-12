package web

import (
	"fmt"
	"net/http"
	"strings"
)

// this only works due to Heroku funneling both http and https requests to the same application port -
// the Heroku router also cleans up request headers, so information in it is assumed reliable
func httpsRedirectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// if the request is not secure or the request is made to the root (or both), redirect to secure canonical name
		if strings.ToLower(req.Header.Get("x-forwarded-proto")) == "http" ||
			req.Host == "kristiandarlington.com" {
			http.Redirect(res, req, fmt.Sprintf("https://www.kristiandarlington.com%s", req.URL), http.StatusPermanentRedirect)
			return
		}

		next.ServeHTTP(res, req)
	})
}
