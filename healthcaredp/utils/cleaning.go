package utils

import (
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"healthcaredp/model"
	"strings"
)

func capitalizeWords(s string) string {
	lower := strings.TrimSpace(strings.ToLower(s))
	words := strings.Fields(lower)
	//create a case with language tag english
	caser := cases.Title(language.English)
	for i, word := range words {
		words[i] = caser.String(word) // Apply title case to each word
	}
	return strings.Join(words, " ")
}

func CleanDataset(scope beam.Scope, col beam.PCollection) beam.PCollection {
	scope = scope.Scope("CleanDataset")
	nameCleaned := beam.ParDo(scope, cleanAdmissionName, col)
	return nameCleaned
}

func cleanAdmissionName(admission model.Admission) model.Admission {
	titleName := capitalizeWords(admission.Name)
	admission.Name = titleName
	return admission
}
