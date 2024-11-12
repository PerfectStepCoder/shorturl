package main

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

// CustomAnalyzer - анализатор, запрещающий использование os.Exit в функции main
var CustomAnalyzer = &analysis.Analyzer{
	Name: "noOsExitInMain",
	Doc:  "check for os.Exit calls in the main function of main package",
	Run:  runCustomAnalyzer,
}

// runCustomAnalyzer - основная логика работы анализа
func runCustomAnalyzer(pass *analysis.Pass) (interface{}, error) {

	for _, file := range pass.Files {
		// Проверяем, что это main пакет
		if pass.Pkg.Name() != "main" {
			continue
		}

		// Проходим по всем узлам AST
		ast.Inspect(file, func(n ast.Node) bool {
			// Ищем функции main
			if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "main" {
				// Ищем os.Exit внутри main
				ast.Inspect(fn.Body, func(n ast.Node) bool {
					if call, ok := n.(*ast.CallExpr); ok {
						if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
							if pkg, ok := selector.X.(*ast.Ident); ok && pkg.Name == "os" && selector.Sel.Name == "Exit" {
								pass.Reportf(call.Pos(), "direct call to os.Exit in main function is prohibited")
							}
						}
					}
					return true
				})
			}
			return true
		})
	}

	return nil, nil
}

func main() {

	// Включаем стандартные анализаторы
	mychecks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	}

	// Включаем все анализаторы класса SA из staticcheck
	for _, v := range staticcheck.Analyzers {
		switch v.Analyzer.Name[:2] {
		case "SA": // анализаторы класса SA
			mychecks = append(mychecks, v.Analyzer)
		case "S": // добавляем по крайней мере один анализатор из класса S
			mychecks = append(mychecks, v.Analyzer)
		case "QF": // добавляем по крайней мере один анализатор из класса QF
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	fmt.Print(mychecks)

	// Добавляем кастомный анализатор
	mychecks = append(mychecks, CustomAnalyzer)

	// Запускаем multichecker
	multichecker.Main(mychecks...)

}
