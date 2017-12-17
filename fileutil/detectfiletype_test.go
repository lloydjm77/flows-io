package fileutil

import (
	"os"
	"strings"
	"testing"
)

func TestDetectFileType(t *testing.T) {
	home := os.Getenv("HOME")
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "should_be_text_file", args: args{path: home + "/software-development/strategy-generator/.gitignore"}, want: "text/plain"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectFileType(tt.args.path); !strings.Contains(got, tt.want) {
				t.Errorf("DetectFileType() = %v, want %v", got, tt.want)
			}
		})
	}
}
