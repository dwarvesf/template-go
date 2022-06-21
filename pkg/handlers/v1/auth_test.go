package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mockEntity "github.com/dwarvesf/go-template/mocks"
	"github.com/dwarvesf/go-template/pkg/config"
	"github.com/dwarvesf/go-template/pkg/consts"
	"github.com/dwarvesf/go-template/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Login(t *testing.T) {
	var cfg = config.GetConfig()

	tests := []struct {
		name       string
		givenBody  string
		givenEmail string
		givenPwd   string
		mockUser   *model.User
		mockErr    error
		wantErr    error
		want       any
		wantCode   int
	}{
		{
			name:      "error - bad request",
			givenBody: `{"email": "", "password": ""}`,
			wantErr:   errors.New("Key: 'LoginRequestBody.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'LoginRequestBody.Password' Error:Field validation for 'Password' failed on the 'required' tag"),
			wantCode:  http.StatusBadRequest,
		},
		{
			name:       "error - unable to login",
			givenBody:  `{"email": "test@test.com", "password": "123456"}`,
			givenEmail: "test@test.com",
			givenPwd:   "123456",
			mockUser:   &model.User{},
			mockErr:    errors.New("invalid email or password"),
			wantErr:    errors.New("invalid email or password"),
			wantCode:   http.StatusUnauthorized,
		},
		{
			name:       "success",
			givenBody:  `{"email": "success@test.com", "password": "123456"}`,
			givenEmail: "success@test.com",
			givenPwd:   "123456",
			mockUser: &model.User{
				Email:       "success@test.com",
				Password:    "123456",
				AccessToken: "token",
			},
			mockErr:  nil,
			wantErr:  nil,
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(tt.givenBody))

			entities := &mockEntity.Entity{}
			entities.On("LoginUser", mock.Anything, tt.givenEmail, tt.givenPwd).Return(tt.mockUser, tt.mockErr)

			h := &Handler{
				cfg:      cfg,
				entities: entities,
			}

			h.Login(ctx)

			var got response

			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))

			if tt.wantErr != nil && tt.wantErr.Error() != got.Error {
				t.Errorf("Handler.Login() code = %v, want %v, error: %s", w.Code, tt.wantCode, got.Error)
				return
			}

			if tt.wantErr == nil {
				require.NotEmpty(t, w.HeaderMap.Get("Set-Cookie"))
				require.NotEmpty(t, got.Data)
			}
		})
	}
}

func TestHandler_Logout(t *testing.T) {
	var cfg = config.GetConfig()

	tests := []struct {
		name string
	}{
		{name: "success"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
			ctx.Request.Header.Add("cookie", fmt.Sprintf("%s=access-token", consts.CookieKey))

			k, _ := ctx.Cookie(consts.CookieKey)
			require.Equal(t, "access-token", k)

			entities := &mockEntity.Entity{}
			h := &Handler{
				cfg:      cfg,
				entities: entities,
			}

			h.Logout(ctx)
			require.Equal(t, http.StatusOK, w.Code)
			var got response

			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
			require.Equal(t, "access_token=; Path=/; Max-Age=0; HttpOnly", w.HeaderMap.Get("Set-Cookie"))
			require.Equal(t, "ok", got.Status)
			require.Equal(t, "logged out", got.Message)
		})
	}
}

func TestHandler_GetLoggedInUser(t *testing.T) {
	var cfg = config.GetConfig()

	tests := []struct {
		name       string
		givenEmail string
		mockUser   *model.User
		mockErr    error
		wantErr    error
		want       any
		wantCode   int
	}{
		{
			name:     "error - unauthorized",
			mockErr:  errors.New("unable to get current user"),
			wantErr:  errors.New("unable to get current user"),
			wantCode: http.StatusUnauthorized,
		},
		{
			name:       "success",
			givenEmail: "success@test.com",
			mockUser: &model.User{
				Email: "success@test.com",
			},
			mockErr:  nil,
			wantErr:  nil,
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest("GET", "/api/v1/auth/me", nil)
			ctx.Set("email", tt.givenEmail)
			entities := &mockEntity.Entity{}
			entities.On("GetUserByEmail", mock.Anything, tt.givenEmail).Return(tt.mockUser, tt.mockErr)

			h := &Handler{
				cfg:      cfg,
				entities: entities,
			}

			h.GetLoggedInUser(ctx)

			var got response

			require.Equal(t, tt.wantCode, w.Code)
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
			if tt.wantErr != nil && tt.wantErr.Error() != got.Error {
				t.Errorf("Handler.GetLoggedInUser() code = %v, want %v, error: %s", w.Code, tt.wantCode, got.Error)
				return
			}

			if tt.wantErr == nil {
				require.NotEmpty(t, got.Data)
			}
		})
	}
}
