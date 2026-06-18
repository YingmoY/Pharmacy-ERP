package app

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"medicare-gateway/internal/config"
	"medicare-gateway/internal/db"
	"medicare-gateway/internal/gateway"
	"medicare-gateway/internal/httpapi"
	"medicare-gateway/internal/medicare"
	"medicare-gateway/internal/queue"
)

type App struct {
	store  *db.Store
	queue  gateway.Queue
	router http.Handler
	cancel context.CancelFunc
}

func New(cfg config.Config, logger *slog.Logger) (*App, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	store, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	if err := store.Migrate(ctx); err != nil {
		store.Close()
		return nil, err
	}

	var q gateway.Queue
	if cfg.EnableRabbitMQ {
		rabbit, err := queue.Connect(cfg.RabbitURL, logger)
		if err != nil {
			logger.Warn("rabbitmq unavailable, async queue disabled", "error", err)
		} else {
			q = rabbit
			logger.Info("rabbitmq async queue enabled")
		}
	}

	client := medicare.NewClient(cfg.MedicareBaseURL, cfg.RequestTimeout)
	service := gateway.NewService(cfg, store, client, q, logger)
	workerCtx, workerCancel := context.WithCancel(context.Background())
	if rabbit, ok := q.(interface {
		Consume(context.Context, func(context.Context, gateway.TaskMessage) error) error
	}); ok {
		if err := rabbit.Consume(workerCtx, service.ProcessTask); err != nil {
			logger.Warn("rabbitmq consumer disabled", "error", err)
		}
	}
	return &App{store: store, queue: q, router: httpapi.New(service), cancel: workerCancel}, nil
}

func (a *App) Router() http.Handler {
	return a.router
}

func (a *App) Close() {
	if a.cancel != nil {
		a.cancel()
	}
	if a.queue != nil {
		_ = a.queue.Close()
	}
	if a.store != nil {
		_ = a.store.Close()
	}
}
