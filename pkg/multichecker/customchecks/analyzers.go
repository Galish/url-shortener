// Package customchecks defines a list of custom static analyzers.
package customchecks

import (
	"golang.org/x/tools/go/analysis"

	"github.com/Galish/url-shortener/pkg/multichecker/customchecks/osexitcheckanalyzer"
)

// Analyzers represents a list of analyzers.
var Analyzers = []*analysis.Analyzer{
	osexitcheckanalyzer.New(),
}
