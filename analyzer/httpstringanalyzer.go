package analyzer

import (
	"regexp"
	"strings"
)

// HTTPStringAnalyzer implementation of StringAnalyzer.
type HTTPStringAnalyzer struct{}

var excludedDomains = []string{"w3", "java", "maven", "xmlsoap", "struts.apache.org", "jakarta.apache.org"}
var regex = regexp.MustCompile("(https?:\\/\\/)([\\da-z\\.-]+\\.[a-z\\.]{2,6})(:[0-9]*)?([\\/\\w \\.-]*)*\\/?")

// Analyze determines if http data is found in a string.
func (hsa HTTPStringAnalyzer) Analyze(s string) AnalysisResult {
	result := regex.FindStringSubmatch(s)

	if result != nil && len(result) > 0 {
		// Since go doesn't support negative lookahead, we are manually excluding results
		// instead of putting them directly in the regex.
		for _, excludedDomain := range excludedDomains {
			if strings.Contains(result[2], excludedDomain) {
				return AnalysisResult{}
			}
		}

		return AnalysisResult{Type: "HTTP", Value: result[0]}
	}

	return AnalysisResult{}
}
