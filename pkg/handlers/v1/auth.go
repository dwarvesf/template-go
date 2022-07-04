package v1

import (
	"net/http"
	"time"

	"github.com/dwarvesf/go-template/pkg/consts"
	"github.com/dwarvesf/go-template/pkg/monitoring"
	"github.com/gin-gonic/gin"
)

type LoginRequestBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login login user with email and pwd this method also set cookie
func (h *Handler) Login(c *gin.Context) {
	m := monitoring.FromContext(c.Request.Context())

	var body LoginRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		m.Errorf(err, "[handler.Login] ShouldBindJSON(&body)")
		c.JSON(http.StatusBadRequest, response{Error: err.Error()})
		return
	}

	resp, err := h.entities.LoginUser(c.Request.Context(), body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response{Error: err.Error()})
		return
	}

	now := time.Now()
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(consts.CookieKey, resp.AccessToken, int(now.Add(time.Duration(h.cfg.AccessTokenTTL)).Unix()), "/", "", true, true)

	c.JSON(http.StatusOK, response{Data: resp})
}

// GetLoggedInUser get logged user via email in request context
func (h *Handler) GetLoggedInUser(c *gin.Context) {
	user, err := h.entities.GetUserByEmail(c.Request.Context(), c.GetString("email"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response{Status: "ok", Data: user})
}

// Logout clear request cookie
func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie(consts.CookieKey, "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, response{
		Status:  "ok",
		Message: "logged out",
	})
}
