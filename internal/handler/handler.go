package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/omekov/dubaicarkzv2/internal/usecase"
)

type Dependencies struct {
	AssetsFS http.FileSystem
	UseCase  usecase.UseCase
}

type hadlerFunc func(w http.ResponseWriter, r *http.Request) error

func RegisterRoutes(r chi.Router, deps Dependencies) {
	home := homeHandler{
		deps.UseCase,
	}
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		// Разрешаем все домены
		AllowedOrigins:   []string{"*"}, // Можете использовать "*"
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Определяет как долго результат запроса может кешироваться (в секундах)
	}))

	AssetsFS := http.FileSystem(deps.AssetsFS)
	fs := http.FileServer(AssetsFS)
	r.Handle("/*", http.StripPrefix("/", fs))
	r.Handle("/static/*", http.StripPrefix("/static", fs))

	// Обработка всех остальных маршрутов Angular через index.html
	r.NotFound(http.FileServer(http.Dir("./dist")).ServeHTTP)
	r.Get("/transport", handler(home.handlerTransport))
	r.Post("/assesstment", handler(home.handlerAssessment))

}

func handler(h hadlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			handlerError(w, r, err)
		}
	}
}

func handlerError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("error during request", slog.String("err", err.Error()))
	http.Error(w, "Something went wrong", http.StatusInternalServerError)
}
