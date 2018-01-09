package fileutil

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestFindSourceFiles(t *testing.T) {
	dir := "../test-files"
	fmt.Print(dir)
	type args struct {
		dir string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "should_find_only_text_files", args: args{dir: dir}, want: []string{
			dir + "/maven-wrapper.properties",
			dir + "/pom.xml",
			dir + "/README.md",
			dir + "/StrategyGeneratorApplication.java",
			dir + "/ClaimWorkflow.wsdl",
			dir + "/application.yml",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := FindSourceFiles(tt.args.dir)
			sort.Strings(files)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(files, tt.want) {
				t.Errorf("FindSourceFiles() = %v, want %v", files, tt.want)
			}
		})
	}
}

func Test_visit(t *testing.T) {
	type args struct {
		files *[]string
	}
	tests := []struct {
		name string
		args args
		want filepath.WalkFunc
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := visit(tt.args.files); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("visit() = %v, want %v", got, tt.want)
			}
		})
	}
}
