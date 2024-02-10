// Package publicchecks defines a list of community static analyzers.
//
// Includes the following list of checks:
//   - `errname` checks that sentinel errors are prefixed with the Err and error types are suffixed with the Error
//   - `bodyclose` checks whether HTTP response body is closed.
package publicchecks

import (
	"github.com/Antonboom/errname/pkg/analyzer"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
)

// Analyzers represents a list of analyzers.
var Analyzers = []*analysis.Analyzer{
	analyzer.New(),
	bodyclose.Analyzer,
}
