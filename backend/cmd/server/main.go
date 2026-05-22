package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/luzyver/tunlify/internal/config"
	"github.com/luzyver/tunlify/internal/db"
	"github.com/luzyver/tunlify/internal/handler"
	"github.com/luzyver/tunlify/internal/middleware"
	"github.com/luzyver/tunlify/internal/service"
	"github.com/luzyver/tunlify/internal/ws"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	database, err := db.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}
	defer database.Close()

	hub := ws.NewHub()
	go hub.Run()

	cfdSvc := service.NewCloudflared(cfg, hub)
	cfdSvc.Logs()
	auditSvc := service.NewAuditLogger(database)
	configSvc := service.NewConfigManager(cfg, database)

	authMw := middleware.NewAuth(cfg.JWTSecret)

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	projectsSvc := service.NewProjects(database, cfg.SelfPath)

	authHandler := handler.NewAuth(cfg, database, auditSvc)
	statusHandler := handler.NewStatus(cfdSvc, cfg)
	controlHandler := handler.NewControl(cfdSvc, auditSvc)
	logsHandler := handler.NewLogs(hub, cfg.JWTSecret)
	configHandler := handler.NewConfig(configSvc, auditSvc)
	tcpHandler := handler.NewTcpAccess(cfdSvc)
	auditHandler := handler.NewAudit(auditSvc)
	healthHandler := handler.NewHealthCheck(cfg.CloudflaredConfig)
	metricsHandler := handler.NewMetrics(cfg)
	settingsHandler := handler.NewSettings(database, auditSvc)
	notifHandler := handler.NewNotifications("./data")
	projectsHandler := handler.NewProjects(projectsSvc, auditSvc)

	r.Post("/auth/login", authHandler.Login)
	r.Get("/api/logs/ws", logsHandler.Stream)

	r.Group(func(r chi.Router) {
		r.Use(authMw.Verify)

		r.Post("/auth/logout", authHandler.Logout)
		r.Post("/auth/refresh", authHandler.Refresh)

		r.Get("/api/status", statusHandler.Get)

		r.Post("/api/control/{action}", controlHandler.Execute)

		r.Get("/api/logs/history", logsHandler.History)

		r.Get("/api/config", configHandler.Get)
		r.Put("/api/config", configHandler.Update)
		r.Get("/api/config/backups", configHandler.ListBackups)
		r.Get("/api/config/backups/{id}", configHandler.GetBackup)
		r.Post("/api/config/validate", configHandler.Validate)

		r.Post("/api/tcp-access/generate", tcpHandler.Generate)
		r.Post("/api/tcp-access/run", tcpHandler.Run)
		r.Delete("/api/tcp-access/run/{pid}", tcpHandler.Stop)

		r.Get("/api/audit", auditHandler.List)

		r.Get("/api/health", healthHandler.Check)
		r.Get("/api/metrics", metricsHandler.Get)
		r.Post("/api/settings/password", settingsHandler.ChangePassword)
		r.Get("/api/notifications", notifHandler.Get)
		r.Put("/api/notifications", notifHandler.Update)
		r.Post("/api/notifications/test", notifHandler.Test)

		r.Get("/api/projects", projectsHandler.List)
		r.Post("/api/projects", projectsHandler.Create)
		r.Put("/api/projects/{id}", projectsHandler.Update)
		r.Delete("/api/projects/{id}", projectsHandler.Delete)
		r.Get("/api/projects/{id}/history", projectsHandler.History)
		r.Get("/api/projects/{id}/output", projectsHandler.Output)
		r.Post("/api/projects/{id}/{action}", projectsHandler.Action)
	})

	srv := &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 600 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("server listening on %s", cfg.ListenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}
