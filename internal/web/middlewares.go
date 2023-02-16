package web

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/nickzhog/gophermart/internal/repositories"
	"github.com/nickzhog/gophermart/internal/web/session"
	"github.com/nickzhog/gophermart/pkg/logging"
)

func SessionMiddleware(logger *logging.Logger, reps repositories.Repositories) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sID, err := session.GetSessionFromCookie(r)
			if err != nil {
				logger.Error(err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			s, err := reps.Session.FindByID(r.Context(), sID)
			if err != nil {
				logger.Error(err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			usr, err := reps.User.FindByID(r.Context(), s.UserID)
			if err != nil {
				logger.Error(err)
				reps.Session.Disable(r.Context(), sID)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			r = session.PutSessionDataInRequest(r, s.ID, usr.ID)
			next.ServeHTTP(w, r)
		})
	}
}

// Gzip compress

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipCompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func gzipDecompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gzReader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer gzReader.Close()
			r.Body = gzReader
		}
		next.ServeHTTP(w, r)
	})
}
