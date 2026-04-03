package transports

import (
	"context"
	"fmt"
	"net/http"
	"sagio/internal/infra/configs"
	"sagio/internal/interfaces/controllers"

	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
)

type HTTPTransport struct {
	conf   *configs.HTTPTransportConfig
	engine *gin.Engine
	server *http.Server
}

func NewHTTPTransport(conf *configs.HTTPTransportConfig) (*HTTPTransport, error) {
	if err := env.Parse(conf); err != nil {
		return nil, err
	}
	return &HTTPTransport{
		conf:   conf,
		engine: gin.Default(),
	}, nil
}

func (t *HTTPTransport) Start(ctx context.Context, resource string) error {
	t.engine.GET(fmt.Sprintf("/%s", resource), controllers.InitTransaction)
	return nil
}
