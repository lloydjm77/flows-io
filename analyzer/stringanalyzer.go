package analyzer

type StringAnalyzer interface {
	Analyze(string) AnalysisResult
}

type AnalysisResult struct {
	Type  string
	Path  string
	Value string
}
