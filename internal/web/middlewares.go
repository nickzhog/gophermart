package web

import (
	"net/http"
)

func (h *HandlerData) logMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		h.Logger.Tracef("request: %s, %s from %s",
			r.Method, r.URL.Path, r.RemoteAddr)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
