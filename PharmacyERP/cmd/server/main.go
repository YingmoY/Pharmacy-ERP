package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/ai"
	"github.com/YingmoY/PharmacyERP/internal/alert"
	"github.com/YingmoY/PharmacyERP/internal/audit"
	auditConsumer "github.com/YingmoY/PharmacyERP/internal/audit/consumer"
	"github.com/YingmoY/PharmacyERP/internal/auth"
	"github.com/YingmoY/PharmacyERP/internal/dashboard"
	"github.com/YingmoY/PharmacyERP/internal/drug"
	"github.com/YingmoY/PharmacyERP/internal/file"
	"github.com/YingmoY/PharmacyERP/internal/inbound"
	"github.com/YingmoY/PharmacyERP/internal/inventory"
	"github.com/YingmoY/PharmacyERP/internal/location"
	"github.com/YingmoY/PharmacyERP/internal/notification"
	"github.com/YingmoY/PharmacyERP/internal/pharmacist"
	"github.com/YingmoY/PharmacyERP/internal/pkg/config"
	"github.com/YingmoY/PharmacyERP/internal/pkg/database"
	"github.com/YingmoY/PharmacyERP/internal/pkg/logger"
	"github.com/YingmoY/PharmacyERP/internal/pkg/medicare"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/YingmoY/PharmacyERP/internal/pkg/mq"
	"github.com/YingmoY/PharmacyERP/internal/report"
	"github.com/YingmoY/PharmacyERP/internal/role"
	"github.com/YingmoY/PharmacyERP/internal/router"
	"github.com/YingmoY/PharmacyERP/internal/sales"
	"github.com/YingmoY/PharmacyERP/internal/shelving"
	"github.com/YingmoY/PharmacyERP/internal/supplier"
	"github.com/YingmoY/PharmacyERP/internal/task"
	"github.com/YingmoY/PharmacyERP/internal/trace"
	"github.com/YingmoY/PharmacyERP/internal/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		panic(fmt.Errorf("load config failed: %w", err))
	}

	log, err := logger.New(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		panic(fmt.Errorf("init logger failed: %w", err))
	}
	defer func() {
		_ = log.Sync()
	}()

	if err := database.Init(cfg.Database, log); err != nil {
		log.Fatal("database init failed", zap.Error(err))
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Error("close database failed", zap.Error(err))
		}
	}()

	mqClient, err := mq.NewClient(cfg.RabbitMQ)
	if err != nil {
		log.Warn("rabbitmq unavailable, operation logs will be skipped", zap.Error(err))
		mqClient = nil
	}
	defer func() {
		if err := mqClient.Close(); err != nil {
			log.Error("close rabbitmq failed", zap.Error(err))
		}
	}()

	// 后台启动操作日志消费者（异步写入审计日志）
	consumerCtx, cancelConsumers := context.WithCancel(context.Background())
	defer cancelConsumers()
	if mqClient != nil {
		logConsumer := auditConsumer.NewLogConsumer(database.DB(), mqClient, log)
		go func() {
			if err := logConsumer.Start(consumerCtx); err != nil {
				log.Error("操作日志消费者异常退出", zap.Error(err))
			}
		}()
	}

	if cfg.Server.Mode != "" {
		gin.SetMode(cfg.Server.Mode)
	}
	engine := gin.New()
	engine.Use(middleware.RequestID(), middleware.Recovery(log), middleware.RequestLogger(log))

	// 初始化医保网关客户端（可选，disabled 时为 nil）
	var medicareClient *medicare.Client
	if cfg.Medicare.Enabled {
		medicareClient = medicare.NewClient(medicare.Config{
			Enabled: cfg.Medicare.Enabled,
			BaseURL: cfg.Medicare.BaseURL,
			Timeout: cfg.Medicare.Timeout,
		})
		log.Info("medicare gateway client initialised", zap.String("base_url", cfg.Medicare.BaseURL))
	}

	// 初始化所有业务模块
	jwtCfg := &cfg.JWT
	secret := cfg.JWT.Secret
	db := database.DB()

	router.RegisterAPIRoutes(engine,
		// 认证与用户管理
		auth.NewModule(db, log, mqClient, jwtCfg),
		user.NewModule(db, log, jwtCfg),
		role.NewModule(db, log, jwtCfg),

		// 基础数据
		drug.NewModule(db, log, mqClient),
		supplier.NewModule(db, log, mqClient),
		location.NewModule(db, log, mqClient),

		// 文件管理
		file.NewModule(db, log, ""),

		// 入库与库存
		inbound.NewModule(db, log, mqClient, secret),
		inventory.NewModule(db, log, secret),
		shelving.NewModule(db, log, secret),

		// 任务管理
		task.NewModule(db, log, mqClient, secret),

		// AI 识别
		ai.NewModule(db, log, mqClient, secret, cfg.AIService.BaseURL, cfg.AIService.Timeout),

		// 销售业务
		sales.NewModule(db, log, mqClient, medicareClient, secret),
		pharmacist.NewModule(db, log, secret),

		// 追溯查询
		trace.NewModule(db, log, secret),

		// 告警与通知
		alert.NewModule(db, log, mqClient, secret),
		notification.NewModule(db, log, mqClient, secret),

		// 仪表盘与报表
		dashboard.NewModule(db, log, mqClient, secret),
		report.NewModule(db, log, mqClient, secret),

		// 审计日志
		audit.NewModule(db, log, secret),
	)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Info("server started", zap.String("addr", srv.Addr), zap.String("env", cfg.App.Env))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server run failed", zap.Error(err))
		}
	}()

	waitForShutdown(srv, cfg.Server.ShutdownTimeout, log, cancelConsumers)
}

func waitForShutdown(srv *http.Server, timeout time.Duration, log *zap.Logger, cancelConsumers context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 先停止 MQ 消费者，再关闭 HTTP 服务
	cancelConsumers()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("graceful shutdown failed", zap.Error(err))
		return
	}
	log.Info("server exited")
}
