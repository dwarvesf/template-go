package routes

import (
	"fmt"
	"testing"

	"github.com/dwarvesf/go-template/pkg/config"
	v1 "github.com/dwarvesf/go-template/pkg/handlers/v1"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestNewRoutes(t *testing.T) {
	cfg := config.GetConfig()
	expectedRoutes := map[string]bool{
		"POST /api/v1/auth/login":  true,
		"POST /api/v1/auth/logout": true,
		"GET /api/v1/auth/me":      true,
	}
	h, err := v1.New(cfg, nil)
	require.NoError(t, err)

	router := NewRoutes(gin.New(), cfg, h)

	routeInfo := router.Routes()

	for _, r := range routeInfo {
		require.NotNil(t, r.HandlerFunc)
		require.Equal(t, true, expectedRoutes[fmt.Sprintf("%s %s", r.Method, r.Path)])
	}
}
