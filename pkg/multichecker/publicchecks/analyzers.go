// Package publicchecks defines a list of static community analyzers.
package publicchecks

import (
	"github.com/Antonboom/errname/pkg/analyzer"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
)

var Analyzers = []*analysis.Analyzer{
	analyzer.New(),     // Checks that sentinel errors are prefixed with the Err and error types are suffixed with the Error
	bodyclose.Analyzer, // Checks whether HTTP response body is closed.
}
