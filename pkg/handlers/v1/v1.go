package v1

import (
	"net/http"

	"github.com/dwarvesf/go-template/pkg/config"
	"github.com/dwarvesf/go-template/pkg/entities"
	"github.com/dwarvesf/go-template/pkg/repo"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg      config.Config
	entities entities.Entity
}

// New will return an instance of Auth struct
func New(cfg config.Config, s repo.Store) (*Handler, error) {
	handler := &Handler{
		cfg:      cfg,
		entities: entities.New(cfg, s),
	}

	return handler, nil
}

// Healthz handler
func (h *Handler) Healthz(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	c.String(http.StatusOK, "OK")
}
