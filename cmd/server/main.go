package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/dwarvesf/go-template/pkg/config"
	"github.com/dwarvesf/go-template/pkg/consts"
	v1 "github.com/dwarvesf/go-template/pkg/handlers/v1"
	"github.com/dwarvesf/go-template/pkg/monitoring"
	"github.com/dwarvesf/go-template/pkg/repo"
	"github.com/dwarvesf/go-template/pkg/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetConfig()

	// init monitor
	monitor, err := monitoring.New(
		cfg,
		cfg.Env != consts.EnvProd,
	)
	if err != nil {
		log.Fatalf("Error in Monitor setup: %s", err)
	}

	defer monitor.Flush(monitoring.DefaultFlushWait)
	ctx := monitoring.SetInContext(context.Background(), monitor)

	s, close := repo.NewPostgresStore(&cfg)
	defer func() {
		err := close()
		if err != nil {
			monitor.Errorf(err, "Error closing postgres store")
		}
	}()

	router := setupRouter(cfg, s)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	shutdownServer(ctx, srv, cfg.GetShutdownTimeout())
}

func shutdownServer(ctx context.Context, srv *http.Server, timeout time.Duration) {
	l := monitoring.FromContext(ctx)
	l.Infof("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		l.Errorf(err, "Server Shutdown:")
	}

}

func setupRouter(cfg config.Config, s repo.Store) *gin.Engine {
	r := gin.New()
	pprof.Register(r)
	r.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/healthz"),
		gin.Recovery(),
	)

	corsOrigins := cfg.GetCORS()

	h, err := v1.New(cfg, s)
	if err != nil {
		log.Fatal(err)
	}

	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowOrigins := corsOrigins

		// allow all localhost and all GET method
		if origin != "" && (strings.Contains(origin, "http://localhost") || c.Request.Method == "GET") {
			allowOrigins = []string{origin}
		} else {
			// support wildcard cors: https://*.domain.com
			for _, url := range allowOrigins {
				if strings.Contains(origin, strings.Replace(url, "https://*", "", 1)) {
					allowOrigins = []string{origin}
					break
				}
			}
		}

		cors.New(
			cors.Config{
				AllowOrigins: allowOrigins,
				AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
				AllowHeaders: []string{"Origin", "Host",
					"Content-Type", "Content-Length",
					"Accept-Encoding", "Accept-Language", "Accept",
					"X-CSRF-Token", "Authorization", "X-Requested-With", "X-Access-Token"},
				ExposeHeaders:    []string{"MeAllowMethodsntent-Length"},
				AllowCredentials: true,
			},
		)(c)
	})

	// handlers
	r.GET("/healthz", h.Healthz)
	routes.NewRoutes(r, cfg, h)

	return r

}
