package web

import "net/http"

func (h *HandlerData) writeError(w http.ResponseWriter, err string, code int) {
	h.Logger.Tracef("write error: err(%s), code(%v)", err, code)
	http.Error(w, err, code)
}

func (h *HandlerData) writeAnswer(w http.ResponseWriter, ans string, code int) {
	h.Logger.Tracef("write answer: ans(%s), code(%v)", ans, code)
	w.WriteHeader(code)
	w.Write([]byte(ans))
}
