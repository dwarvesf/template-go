package entities

import (
	"context"
	"reflect"
	"testing"

	mockRepo "github.com/dwarvesf/go-template/mocks"
	mockUserStore "github.com/dwarvesf/go-template/mocks/user"
	"github.com/dwarvesf/go-template/pkg/config"
	"github.com/dwarvesf/go-template/pkg/model"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestEntity_GetUserByEmail(t *testing.T) {
	var cfg = config.GetConfig()

	tests := []struct {
		name       string
		givenEmail string
		mockRsp    *model.User
		mockErr    error
		want       *model.User
		wantErr    bool
	}{
		{
			name:       "success",
			givenEmail: "test@sample.com",
			mockRsp:    &model.User{Email: "test@sample.com", Password: "123456"},
			mockErr:    nil,
			want:       &model.User{Email: "test@sample.com", Password: "123456"},
			wantErr:    false,
		},
		{
			name:       "error",
			givenEmail: "notfound@sample.com",
			mockRsp:    nil,
			mockErr:    gorm.ErrRecordNotFound,
			want:       nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRp := &mockUserStore.Store{}
			userRp.On("GetUserByEmail", mock.Anything, tt.givenEmail).Return(tt.mockRsp, tt.mockErr)

			rp := &mockRepo.Store{}
			rp.On("UserRepo").Return(userRp)

			e := &entity{
				cfg:  cfg,
				repo: rp,
			}

			got, err := e.GetUserByEmail(context.Background(), tt.givenEmail)
			if (err != nil) != tt.wantErr {
				t.Errorf("Entity.GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Entity.GetUserByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntity_CreateUser(t *testing.T) {
	var cfg = config.GetConfig()

	tests := []struct {
		name      string
		givenUser model.User
		mockRsp   *model.User
		mockErr   error
		want      *model.User
		wantErr   bool
	}{
		{
			name:      "success",
			givenUser: model.User{Email: "test@sample.com", Password: "123456"},
			mockRsp:   &model.User{Email: "test@sample.com", Password: "123456"},
			mockErr:   nil,
			want:      &model.User{Email: "test@sample.com", Password: "123456"},
			wantErr:   false,
		},
		{
			name:      "error",
			givenUser: model.User{Email: "test@sample.com", Password: "123456"},
			mockRsp:   nil,
			mockErr:   gorm.ErrRecordNotFound,
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRp := &mockUserStore.Store{}
			userRp.On("CreateUser", mock.Anything, mock.Anything).Return(tt.mockRsp, tt.mockErr)

			rp := &mockRepo.Store{}
			rp.On("UserRepo").Return(userRp)

			e := &entity{
				cfg:  cfg,
				repo: rp,
			}

			got, err := e.CreateUser(context.Background(), tt.givenUser)
			if (err != nil) != tt.wantErr {
				t.Errorf("Entity.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Entity.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
