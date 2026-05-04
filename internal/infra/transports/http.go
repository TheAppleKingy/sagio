package transports

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"sagio/internal/infra/configs"

	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
)

type HTTPTransport struct {
	*BaseTransport
	conf   *configs.HTTPTransportConfig
	engine *gin.Engine
	server *http.Server
}

func NewHTTPTransport(
	core *SagioCore,
) (Transport, error) {
	conf := &configs.HTTPTransportConfig{}
	if err := env.Parse(conf); err != nil {
		return nil, err
	}
	gin.SetMode(gin.ReleaseMode)
	res := &HTTPTransport{
		conf:   conf,
		engine: gin.Default(),
	}
	res.BaseTransport = &BaseTransport{
		core:            core,
		Type:            HTTP,
		shutdownTimeout: conf.ShutdownTimeout,
		impl:            res,
		ctx:             core.ctx,
	}
	return res, nil
}

func (t *HTTPTransport) sendResponse(status int, info string, ctx *gin.Context) {
	ctx.JSON(status, gin.H{
		"detail": info,
	})
}

func (t *HTTPTransport) connect() error {
	uri := fmt.Sprintf("/%s", t.core.resource)
	t.engine.GET(uri, func(ctx *gin.Context) {
		select {
		case <-t.ctx.Done():
			slog.Warn("HTTP transport got cancel command. Unable to continue requests handling")
			t.sendResponse(http.StatusServiceUnavailable, "Service is shutdowning now", ctx)
			return
		default:
		}
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			t.sendResponse(http.StatusUnprocessableEntity, err.Error(), ctx)
			return
		}
		key := ctx.GetHeader(configs.AppConfig.StepKeyName)
		if key == "" {
			t.sendResponse(http.StatusBadRequest, "Undefined step type", ctx)
			return
		}
		if err := t.core.StepProcessing(body, key, &t.wg); err != nil {
			t.sendResponse(http.StatusBadRequest, err.Error(), ctx)
			return
		}
		t.sendResponse(http.StatusOK, "ok", ctx)
	})
	t.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", t.conf.Port),
		Handler: t.engine,
	}
	go func() {
		if err := t.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP transport cannot start consuming requests", "error", err.Error(), "uri", uri)
		}
	}()
	return nil
}

func (t *HTTPTransport) disconnect() {
	//nolint:errcheck
	t.server.Shutdown(context.Background())
}
