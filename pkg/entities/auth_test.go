package entities

import (
	"context"
	"errors"
	"testing"

	mockRepo "github.com/dwarvesf/go-template/mocks"
	mockUserStore "github.com/dwarvesf/go-template/mocks/user"
	"github.com/dwarvesf/go-template/pkg/config"
	"github.com/dwarvesf/go-template/pkg/model"
	"github.com/dwarvesf/go-template/pkg/util"
	"github.com/golang-jwt/jwt"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_entity_LoginUser(t *testing.T) {

	var cfg = config.GetConfig()

	type args struct {
		email string
		pwd   string
	}
	checkPasswordHashFunc = func(a, b string) bool {
		return a == b
	}

	issueAccessTokenFunc = func(cfg config.Config, user *model.User) (string, error) {
		if user.Email == "invalid@email.com" {
			return "", errors.New("internal server error")
		}

		return issueAccessToken(cfg, user)
	}

	defer func() {
		checkPasswordHashFunc = util.CheckPasswordHash
		issueAccessTokenFunc = issueAccessToken
	}()

	tests := []struct {
		name     string
		args     args
		mockResp *model.User
		mockErr  error
		wantErr  bool
	}{
		{
			name:     "error not found",
			args:     args{email: "invalid@email.com", pwd: "1234"},
			mockResp: nil,
			mockErr:  errInvalidEmailOrPwd,
			wantErr:  true,
		},
		{
			name:     "error compare pwd",
			args:     args{email: "valid_email@email.com", pwd: "wrong"},
			mockResp: &model.User{Email: "valid_email@email.com", Password: "1234"},
			mockErr:  nil,
			wantErr:  true,
		},
		{
			name:     "error issue token",
			args:     args{email: "invalid@email.com", pwd: "1234"},
			mockResp: &model.User{Email: "invalid@email.com", Password: "1234"},
			mockErr:  nil,
			wantErr:  true,
		},
		{
			name:     "success",
			args:     args{email: "valid_email@email.com", pwd: "1234"},
			mockResp: &model.User{Email: "valid_email@email.com", Password: "1234"},
			mockErr:  nil,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRp := &mockUserStore.Store{}
			userRp.On("GetUserByEmail", mock.Anything, tt.args.email).Return(tt.mockResp, tt.mockErr)

			rp := &mockRepo.Store{}
			rp.On("UserRepo").Return(userRp)

			e := &entity{
				cfg:  cfg,
				repo: rp,
			}

			got, err := e.LoginUser(context.Background(), tt.args.email, tt.args.pwd)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NotNil(t, got)
			require.NotEmpty(t, got.AccessToken)

			payload := jwt.MapClaims{}
			_, err = jwt.ParseWithClaims(got.AccessToken, &payload, func(token *jwt.Token) (interface{}, error) {
				return cfg.JWTSecret, nil
			})
			require.NoError(t, err)
			require.Equal(t, payload["email"], tt.args.email)

		})
	}
}
