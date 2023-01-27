package web

import "net/http"

func showError(w http.ResponseWriter, err string, code int) {
	w.WriteHeader(code)
	w.Write([]byte(err))
}
