package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
)

func TestStaticLint(t *testing.T) {
	// Включаем стандартные анализаторы
	mychecks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	}

	addNewAnalizers(mychecks)

	// Добавляем кастомный анализатор
	mychecks = append(mychecks, CustomAnalyzer)

	assert.Equal(t, len(mychecks), 4)

}
