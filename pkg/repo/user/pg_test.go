package user

import (
	"context"
	"reflect"
	"testing"

	"github.com/dwarvesf/go-template/pkg/config"
	"github.com/dwarvesf/go-template/pkg/model"
	"github.com/dwarvesf/go-template/pkg/util/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var cfg = config.GetConfig()

func TestNewUserRepoPg(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want Store
	}{
		{
			name: "success",
			args: args{db: nil},
			want: &userRepo{db: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserRepoPg(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserRepoPg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepo_GetUserByID(t *testing.T) {
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "success",
			args: args{email: "dwarvesf@dwarvesv.com"},
			want: &model.User{
				ID:       uuid.NullUUID{UUID: uuid.MustParse(("ae342f45-33a0-4d0b-8c03-4b653e401a7c")), Valid: true},
				Email:    "dwarvesf@dwarvesv.com",
				Password: "123456",
			},
		},
		{
			name:    "error not found",
			args:    args{email: "not found"},
			want:    &model.User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutil.WithTxDB(t, cfg, func(tx *gorm.DB) {
				testutil.LoadTestSQLFile(t, tx, "testdata/user.sql")

				s := userRepo{
					db: tx,
				}

				got, err := s.GetUserByEmail(tt.args.ctx, tt.args.email)
				if (err != nil) != tt.wantErr {
					t.Errorf("userRepo.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				ignoreFields := cmpopts.IgnoreFields(
					model.User{},
					"CreatedAt",
					"UpdatedAt",
				)

				if !cmp.Equal(got, tt.want, ignoreFields) {
					t.Errorf("userRepo.GetUserByID() = %v, want %v \n diff: %v", got, tt.want, cmp.Diff(got, tt.want, ignoreFields))
					t.FailNow()
				}
			})
		})

	}
}

func Test_userRepo_CreateUser(t *testing.T) {

	type args struct {
		ctx  context.Context
		user model.User
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "success",
			args: args{user: model.User{Email: "test@dwarvesv.com", Password: "123"}},
			want: &model.User{ID: uuid.NullUUID{Valid: true}, Email: "test@dwarvesv.com", Password: "123"},
		},
		{
			name:    "error email existed",
			args:    args{user: model.User{Email: "dwarvesf@dwarvesv.com", Password: "123"}},
			want:    &model.User{Email: "dwarvesf@dwarvesv.com", Password: "123"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testutil.WithTxDB(t, cfg, func(tx *gorm.DB) {
				testutil.LoadTestSQLFile(t, tx, "testdata/user.sql")

				s := userRepo{
					db: tx,
				}

				got, err := s.CreateUser(tt.args.ctx, tt.args.user)
				if (err != nil) != tt.wantErr {
					t.Errorf("userRepo.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				ignoreFields := cmpopts.IgnoreFields(
					model.User{},
					"ID",
					"CreatedAt",
					"UpdatedAt",
				)

				if !cmp.Equal(got, tt.want, ignoreFields) {
					t.Errorf("userRepo.GetUserByID() = %v, want %v \n diff: %v", got, tt.want, cmp.Diff(got, tt.want, ignoreFields))
					t.FailNow()
				}
			})
		})
	}
}
