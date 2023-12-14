package config

import (
	"os"
	"reflect"
	"testing"

	"bou.ke/monkey"
)

func TestReadConf(t *testing.T) {
	correct_conf := "license_dir: /var/run/licMgr"
	wrong_conf := "wrong conf"

	type args struct {
		file string
	}
	tests := []struct {
		name       string
		args       args
		want       *LicenseConfig
		wantErr    bool
		patchFunc  interface{}
		targetFunc interface{}
	}{
		{
			name:       "Err1: file not exist",
			args:       args{file: "config.conf"},
			want:       nil,
			wantErr:    true,
			patchFunc:  nil,
			targetFunc: nil,
		},
		{
			name:       "Err2: not yaml format",
			args:       args{file: "config.conf"},
			want:       nil,
			wantErr:    true,
			targetFunc: os.ReadFile,
			patchFunc: func(string) ([]byte, error) {
				return []byte(wrong_conf), nil
			},
		},
		{
			name:       "OK",
			args:       args{file: "config.conf"},
			want:       &LicenseConfig{LisDir: "/var/run/licMgr"},
			wantErr:    false,
			targetFunc: os.ReadFile,
			patchFunc: func(string) ([]byte, error) {
				return []byte(correct_conf), nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.patchFunc != nil {
				monkey.Patch(tt.targetFunc, tt.patchFunc)
			}

			got, err := ReadConf(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConf() = %v, want %v", got, tt.want)
			}
		})
	}
}
