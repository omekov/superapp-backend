package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/omekov/dubaicarkzv2/internal/config"
	"github.com/omekov/dubaicarkzv2/internal/handler"
	"github.com/omekov/dubaicarkzv2/internal/usecase"
	"github.com/omekov/dubaicarkzv2/internal/usecase/repository"
)

func Run(ctx context.Context) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", cfg.SqlitePath)
	if err != nil {
		return err
	}
	defer db.Close()

	func(db *sql.DB) {
		var kgdURL string
		_, err := db.Exec(`CREATE TABLE IF NOT EXISTS kgd_data_migration (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			kgd_url TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime'))
		);`)
		if err != nil {
			slog.Error("Exec", err)
			return
		}

		err = db.QueryRow("SELECT kgd_url FROM kgd_data_migration WHERE kgd_url = ? ORDER BY created_at DESC;", cfg.KGDURL).Scan(&kgdURL)
		if errors.Is(err, sql.ErrNoRows) {
			err := startMigrate(db, cfg.KGDURL)
			if err != nil {
				slog.Error("StartMigrate", err)
				return
			}
			_, err = db.Exec("INSERT INTO kgd_data_migration (kgd_url) VALUES  (?)", cfg.KGDURL)
			if err != nil {
				slog.Error("Exec", err)
				return
			}
		} else {
			slog.Error("QueryRow", err)
			return
		}
	}(db)

	repo, err := repository.NewRepository(db)
	if err != nil {
		return fmt.Errorf("repository -> %v", err)
	}
	uc := usecase.NewUseCase(repo)

	r := chi.NewRouter()

	handler.RegisterRoutes(r, handler.Dependencies{
		AssetsFS: http.Dir(cfg.AssetsDir),
		UseCase:  uc,
	})

	s := http.Server{
		Addr:    cfg.ServerAddr,
		Handler: r,
	}

	go func() {
		tb := NewTelegramBot(cfg.TelegramApiToken)
		if err := tb.Init(); err != nil {
			slog.Error("tb.Init", err)
		}
	}()

	go func() {
		<-ctx.Done()
		slog.Info("shutting down server")
		s.Shutdown(ctx)
	}()

	slog.Info("starting server", slog.String("addr", cfg.ServerAddr))
	if err := s.ListenAndServe(); err != http.ErrServerClosed && err != nil {
		return err
	}
	return nil
}
