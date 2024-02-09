// Package multichecker implements a static analysis tool.
package multichecker

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

// New builds an analysis tool containing the specified analyzers.
func New(opts ...func(*[]*analysis.Analyzer)) {
	var checks []*analysis.Analyzer

	for _, opt := range opts {
		opt(&checks)
	}

	multichecker.Main(checks...)
}

// WithAnalyzers adds the specified analyzers to the initial list.
func WithAnalyzers(a []*analysis.Analyzer) func(c *[]*analysis.Analyzer) {
	return func(c *[]*analysis.Analyzer) {
		*c = append(*c, a...)
	}
}
