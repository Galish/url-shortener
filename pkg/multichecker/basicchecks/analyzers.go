// Package basicchecks defines a list of standard static analyzers.
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
