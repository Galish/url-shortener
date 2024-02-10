// Package customchecks defines a list of custom static analyzers.
//
// Includes the following list of checks:
//   - `osexitcheckanalyzer` checks for a direct os.Exit call in the main function of the main package.
package customchecks

import (
	"golang.org/x/tools/go/analysis"

	"github.com/Galish/url-shortener/pkg/multichecker/customchecks/osexitcheckanalyzer"
)

// Analyzers represents a list of analyzers.
var Analyzers = []*analysis.Analyzer{
	osexitcheckanalyzer.New(),
}
