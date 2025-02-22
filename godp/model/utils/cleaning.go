package model

import (
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

// CleanDataset applies cleaning operations on the input dataset to the pipeline.ì
//
// Parameters:
//   - scope: A beam.Scope object representing the current Apache Beam pipeline scope.
//   - col: A beam.PCollection containing the input dataset to be cleaned.
//
// Returns:
//
//	A beam.PCollection containing the cleaned dataset with standardized admission names.
func CleanDataset(scope beam.Scope, col beam.PCollection) beam.PCollection {
	scope = scope.Scope("CleanDataset")
	nameCleaned := beam.ParDo(scope, cleanAdmissionName, col)
	return nameCleaned
}

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

func cleanAdmissionName(admission Admission) Admission {
	titleName := capitalizeWords(admission.Name)
	admission.Name = titleName
	return admission
}
