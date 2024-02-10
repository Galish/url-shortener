// Package main is the entry point to the static analysis tool.
// This package makes it easy to build an analysis tool containing specified analyzers.
package main

import (
	"github.com/Galish/url-shortener/pkg/multichecker"
	"github.com/Galish/url-shortener/pkg/multichecker/basicchecks"
	"github.com/Galish/url-shortener/pkg/multichecker/customchecks"
	"github.com/Galish/url-shortener/pkg/multichecker/publicchecks"
	"github.com/Galish/url-shortener/pkg/multichecker/staticcheck"
)

func main() {
	multichecker.New(
		multichecker.WithAnalyzers(basicchecks.Analyzers),
		multichecker.WithAnalyzers(staticcheck.Analyzers),
		multichecker.WithAnalyzers(publicchecks.Analyzers),
		multichecker.WithAnalyzers(customchecks.Analyzers),
	)
}
