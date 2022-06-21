package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestENVReader_Load(t *testing.T) {
	type args struct {
		v viper.Viper
	}
	tests := []struct {
		name    string
		r       *ENVReader
		args    args
		wantErr bool
	}{
		{
			name:    "success",
			r:       &ENVReader{},
			args:    args{v: *viper.New()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ENVReader{}
			_, err := r.Load(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ENVReader.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
