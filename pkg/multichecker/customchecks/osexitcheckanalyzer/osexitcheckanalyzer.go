// Package osexitcheckanalyzer implement a custom analyzer that checks
// for a direct os.Exit call in the main function of the main package.
package osexitcheckanalyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const message = "illegal `os.Exit` call"

// New returns a new `osexitinmaincheck` analyzer.
func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "osexitinmaincheck",
		Doc:  "check for direct os.Exit call in the main function of the main package",
		Run:  run,
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if !isMainPackage(file) {
			continue
		}

		var parents []*ast.Node

		ast.Inspect(file, func(node ast.Node) bool {
			if node == nil {
				parents = parents[:len(parents)-1]
				return true
			}

			parents = append(parents, &node)

			c, ok := node.(*ast.CallExpr)
			if ok && isOSExitCall(c) && isMainFunc(parents) {
				pass.Reportf(c.Fun.Pos(), message)
			}

			return true
		})
	}

	return nil, nil
}

func isMainPackage(file *ast.File) bool {
	return file.Name.String() == "main"
}

func isMainFunc(chain []*ast.Node) bool {
	for i := len(chain) - 1; i >= 0; i-- {
		node := *chain[i]

		s, ok := node.(*ast.FuncDecl)
		if ok && s.Name.String() == "main" {
			return true
		}
	}

	return false
}

func isOSExitCall(call *ast.CallExpr) bool {
	s, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || s.Sel.Name != "Exit" {
		return false
	}

	if x, ok := s.X.(*ast.Ident); !ok || x.Name != "os" {
		return false
	}

	return true
}
