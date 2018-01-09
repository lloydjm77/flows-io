package fileutil

import (
	"strings"
	"testing"
)

func TestDetectFileType(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "should_be_xml_file", args: args{path: "../test-files/pom.xml"}, want: "text/xml"},
		{name: "should_be_java_text_file", args: args{path: "../test-files/StrategyGeneratorApplication.java"}, want: "text/plain"},
		{name: "should_generate_error", args: args{path: "../test-files/does-not-exist.txt"}, want: "invalid argument"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectFileType(tt.args.path); !strings.Contains(got, tt.want) {
				t.Errorf("DetectFileType() = %v, want %v", got, tt.want)
			}
		})
	}
}
