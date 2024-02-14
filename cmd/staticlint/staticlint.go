// Package main is the entry point to the static analysis tool.
// This package makes it easy to build an analysis tool containing specified analyzers.
//
// # Configuration
//
// Multichecker can be configured by providing a list of desired checker packages:
//   - `basicchecks` provides a set of standard static analyzers.
//   - `staticcheck` provides a set of `staticcheck.io` static analyzers.
//   - `publicchecks` provides a set of community static analyzers.
//   - `customchecks` provides a set of of custom static analyzers.
//
// Here's a code example:
//
//	multichecker.New(
//		multichecker.WithAnalyzers(basicchecks.Analyzers),
//		multichecker.WithAnalyzers(staticcheck.Analyzers),
//		multichecker.WithAnalyzers(publicchecks.Analyzers),
//		multichecker.WithAnalyzers(customchecks.Analyzers),
//	)
//
// # Usage
//
// You can perform static analysis by running the source file:
//
//	go run cmd/staticlint/staticlint.go ./...
//
// Or by running the build and executing the binary:
//
//	go build cmd/staticlint/*.go
//	./staticlint ./...
//
// # Flags
//
// By default all analyzers are run.
// To select specific analyzers, use the -NAME flag for each one, or -NAME=false to run all analyzers not explicitly disabled.
//
// To get a complete list of flags, run the `help` command:
//
//	./staticlint help
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
