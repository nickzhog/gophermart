package web

import "net/http"

func showError(w http.ResponseWriter, err string, code int) {
	http.Error(w, err, code)
}
