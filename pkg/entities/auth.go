package entities

import (
	"context"
	"fmt"
	"time"

	"github.com/dwarvesf/go-template/pkg/config"
	"github.com/dwarvesf/go-template/pkg/model"
	"github.com/dwarvesf/go-template/pkg/monitoring"
	"github.com/dwarvesf/go-template/pkg/util"
	"github.com/golang-jwt/jwt"
)

var checkPasswordHashFunc = util.CheckPasswordHash
var issueAccessTokenFunc = issueAccessToken

func issueAccessToken(cfg config.Config, user *model.User) (string, error) {
	issuer := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"exp":   time.Now().Add(cfg.AccessTokenTTL).Unix(),
		"email": user.Email,
		"uuid":  user.ID.UUID,
	})

	return issuer.SignedString(cfg.JWTSecret)
}

func (e *entity) LoginUser(ctx context.Context, email, pwd string) (*model.User, error) {
	m := monitoring.FromContext(ctx)
	user, err := e.repo.UserRepo().GetUserByEmail(ctx, email)
	if err != nil {
		m.Errorf(err, "[entity.LoginUser] GetUserByEmail(ctx, email=%v)", email)
		return nil, errInvalidEmailOrPwd
	}

	if !checkPasswordHashFunc(pwd, user.Password) {
		return nil, errInvalidEmailOrPwd
	}

	user.AccessToken, err = issueAccessTokenFunc(e.cfg, user)
	if err != nil {
		m.Errorf(err, "[entity.LoginUser] issueAccessTokenFunc(e.cfg, user=%#v)", user)
		return nil, fmt.Errorf("unable to issue new token for email %s with err %s", user.Email, err.Error())
	}

	return user, nil
}
