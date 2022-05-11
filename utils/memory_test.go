package utils

import (
	"testing"
)

func TestGetObjSize(t *testing.T) {
	var strptr *string
	tmp := "hello world"
	strptr = &tmp
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{"case nil", args{nil}, 0, true},
		{"case string empty", args{""}, 4, false}, // 空字符串返回为4
		{"case string len 1", args{"1"}, 5, false},
		{"case string len 2", args{"12"}, 6, false},
		{"case string ptr", args{strptr}, 15, false},
		{"case int", args{1}, 4, false},
		{"case uint", args{uint(1)}, 4, false},
		{"case int32", args{int32(1)}, 4, false},
		{"case int64", args{int64(1)}, 4, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetObjSize(tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetObjSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetObjSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseMemorySize(t *testing.T) {
	type args struct {
		size string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{"case 1k", args{"1k"}, 1, false},
		{"case 1M", args{"1M"}, 1024, false},
		{"case 1mb", args{"1mb"}, 1024, false},
		{"case 1GB", args{"1GB"}, 1048576, false},
		{"case err", args{"1GB1"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMemorySize(tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMemorySize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if int(got) != tt.want {
				t.Errorf("ParseMemorySize() = %v, want %v", got, tt.want)
			}
		})
	}
}
