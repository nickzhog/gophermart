package web

import "net/http"

func writeError(w http.ResponseWriter, err string, code int) {
	http.Error(w, err, code)
}

func writeAnswer(w http.ResponseWriter, ans string, code int) {
	w.WriteHeader(code)
	w.Write([]byte(ans))
}
