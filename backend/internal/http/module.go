package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smoosex/lumina/backend/internal/config"
	"github.com/smoosex/lumina/backend/internal/notes"
	"github.com/smoosex/lumina/backend/internal/uploadauth"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"http",
	fx.Provide(NewRouter, NewServer),
	fx.Invoke(RegisterLifecycle),
)

func NewRouter(
	cfg config.Config,
	noteHandler *notes.Handler,
	auth *uploadauth.Middleware,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	router.Static("/covers", "data/covers")

	api := router.Group("/api")
	noteHandler.RegisterPublic(api)

	admin := api.Group("/admin")
	admin.Use(auth.Handle())
	noteHandler.RegisterAdmin(admin)

	return router
}

func NewServer(cfg config.Config, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: router,
	}
}

func RegisterLifecycle(lc fx.Lifecycle, server *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil &&
					!errors.Is(err, http.ErrServerClosed) {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}
