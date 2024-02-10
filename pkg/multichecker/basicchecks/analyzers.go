// Package basicchecks defines a list of standard static analyzers.
//
// Includes the following list of checks:
//   - `shadow` checks for shadowed variables.
//   - `structtag` checks struct field tags are well formed.
//   - `lostcancel` checks for failure to call a context cancellation function.
//   - `httpresponse` checks for mistakes using HTTP responses.
//   - `loopclosure` checks for references to enclosing loop variables from within nested functions.
//   - `tests` checks for common mistaken usages of tests and examples.
//   - `unmarshal` checks for passing non-pointer or non-interface types to unmarshal and decode functions.
//   - `unreachable` checks for unreachable code.
//   - `unusedresult` checks for unused results of calls to certain pure functions.
package basicchecks

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
)

// Analyzers represents a list of analyzers.
var Analyzers = []*analysis.Analyzer{
	shadow.Analyzer,
	structtag.Analyzer,
	lostcancel.Analyzer,
	httpresponse.Analyzer,
	loopclosure.Analyzer,
	tests.Analyzer,
	unmarshal.Analyzer,
	unreachable.Analyzer,
	unusedresult.Analyzer,
}
