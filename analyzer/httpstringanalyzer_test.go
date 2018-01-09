package analyzer

import (
	"reflect"
	"testing"
)

func TestHTTPStringAnalyzer_Analyze(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		hsa  HTTPStringAnalyzer
		args args
		want AnalysisResult
	}{
		{name: "should_return_url", args: args{s: "http://google.com"}, want: AnalysisResult{Type: "HTTP", Value: "http://google.com"}},
		{name: "should_return_url2", args: args{s: "http://testing.com:8080"}, want: AnalysisResult{Type: "HTTP", Value: "http://testing.com:8080"}},
		{name: "should_exclude_url", args: args{s: "http://maven.com"}, want: AnalysisResult{}},
		{name: "should_exclude_url2", args: args{s: "http://www.w3.com"}, want: AnalysisResult{}},
		{name: "should_find_no_match", args: args{s: "no-match"}, want: AnalysisResult{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hsa := HTTPStringAnalyzer{}
			if got := hsa.Analyze(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HTTPStringAnalyzer.Analyze() = %v, want %v", got, tt.want)
			}
		})
	}
}
