package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestFileReader_Load(t *testing.T) {
	type fields struct {
		filename string
		dirname  string
	}
	type args struct {
		v viper.Viper
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "success",
			fields:  fields{filename: ".env.sample", dirname: "../.."},
			args:    args{v: *viper.New()},
			wantErr: false,
		},
		{
			name:    "error",
			fields:  fields{filename: ".env.sample", dirname: "."},
			args:    args{v: *viper.New()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FileReader{
				filename: tt.fields.filename,
				dirname:  tt.fields.dirname,
			}
			_, err := r.Load(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileReader.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
