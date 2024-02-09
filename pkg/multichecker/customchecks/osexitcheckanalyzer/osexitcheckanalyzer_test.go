package osexitcheckanalyzer_test

import (
	"testing"

	"github.com/Galish/url-shortener/pkg/multichecker/customchecks/osexitcheckanalyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestMyAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), osexitcheckanalyzer.New(), "./...")
}
