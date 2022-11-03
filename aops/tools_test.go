package aops

import (
	"reflect"
	"testing"
)

func Test_getParamsFromID(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "normal test",
			args:    args{id: `@middleware-injection(path:"xxxx")`},
			wantErr: false,
			want: []string{
				`path := "xxxx"`,
			},
		},
		{
			name:    "two params test",
			args:    args{id: `@middleware-injection(path:"xxxx", name:33)`},
			wantErr: false,
			want: []string{
				`path := "xxxx"`,
				`name := 33`,
			},
		},
	}
	ij := injectDetail{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ij.getParamsFromID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("getParamsFromID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getParamsFromID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
