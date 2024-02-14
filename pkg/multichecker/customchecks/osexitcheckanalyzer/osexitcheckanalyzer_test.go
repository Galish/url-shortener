package osexitcheckanalyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/Galish/url-shortener/pkg/multichecker/customchecks/osexitcheckanalyzer"
)

func TestMyAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), osexitcheckanalyzer.New(), "./...")
}
