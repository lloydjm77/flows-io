package analyzer

// StringAnalyzer defines the methods to be implemented by all StringAnalyzer
// implementations.
type StringAnalyzer interface {

	// Analyze reads a string and determines if it contains matches
	// based on the specific implementation.  It returns an AnalysisResult.
	Analyze(string) AnalysisResult
}

// AnalysisResult encapsulates information about the results of an analysis.
type AnalysisResult struct {
	Type  string
	Path  string
	Value string
}
