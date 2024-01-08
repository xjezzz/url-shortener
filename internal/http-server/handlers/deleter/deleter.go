package deleter

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"
)

type URLDeleter interface {
	DeleteURL(alias string) error
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(slog.String("request_id", middleware.GetReqID(r.Context())))
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias in empty")
			render.JSON(w, r, response.Error("not found"))
			return
		}
		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			render.JSON(w, r, response.Error("not found"))

		}
		if err != nil {
			log.Error("failed to deleter the url", sl.Err(err))
			render.JSON(w, r, response.Error("internal error"))
			return
		}
		log.Info("url is deleted", slog.String("alias", alias))
	}

}
