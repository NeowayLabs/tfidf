package go_tfidf_test

import (
	"testing"

	"github.com/NeowayLabs/go-tfidf/similarity"

	go_tfidf "github.com/NeowayLabs/go-tfidf"
	"github.com/stretchr/testify/assert"
)

func TestSetAndGetDocumentSeparator(t *testing.T) {
	ti := go_tfidf.New()
	input := "-"

	ti.SetDocumentSeparator(input)

	assert.Equal(t, input, ti.GetDocumentSeparator())

}

func TestNormalizedTermFrequency(t *testing.T) {
	inputTerms := []string{"valid", "document"}
	expected := map[string]float64{
		"valid":    0.5,
		"document": 0.5,
	}

	actual := go_tfidf.NormalizedTermFrequency(inputTerms)

	assert.Equal(t, expected, actual)
}

func TestInverseDocumentFrequency(t *testing.T) {
	inputTerm := "life"
	inputDocuments := []string{
		"The game of life is a game of everlasting learning",
		"The unexamined life is not worth living",
		"Never stop learning",
	}
	separator := " "
	expected := 1.4054651081081644

	actual := go_tfidf.InverseDocumentFrequency(inputTerm, inputDocuments, separator)

	assert.Equal(t, expected, actual)
}

func TestAddDocumentsWhenInputIsEmpty(t *testing.T) {
	inputDocuments := []string{}
	ti := go_tfidf.New()
	err := ti.AddDocuments(inputDocuments)

	assert.NotNil(t, err)
}

func TestAddDocumentsWhenAtLeastOneDocumentIsInvalid(t *testing.T) {
	inputDocuments := []string{"valid document", ""}
	expectedDocument := make([]string, 0)
	ti := go_tfidf.New()
	err := ti.AddDocuments(inputDocuments)

	assert.Equal(t, expectedDocument, ti.GetDocuments())
	assert.NotNil(t, err)
}

func TestAddDocuments(t *testing.T) {
	inputDocuments := []string{"valid document"}

	expectedDocuments := []string{"valid document"}
	expectedDocumentsTerms := []string{"valid", "document"}
	expectedDocumentsNormTermFrequency := []map[string]float64{
		map[string]float64{
			"valid":    0.5,
			"document": 0.5,
		},
	}

	ti := go_tfidf.New()
	err := ti.AddDocuments(inputDocuments)

	assert.Equal(t, expectedDocuments, ti.GetDocuments())
	assert.Equal(t, expectedDocumentsTerms, ti.GetDocumentsTerms())
	assert.Equal(t, expectedDocumentsNormTermFrequency, ti.GetDocumentsNormTermFrequency())
	assert.Equal(t, expectedDocuments, ti.GetDocuments())
	assert.Nil(t, err)
}

func TestCalculateQueryTermsTfIdfForEachDocumentWhenQueryIsInvalid(t *testing.T) {
	inputDocuments := []string{
		"The game of life is a game of everlasting learning",
		"The unexamined life is not worth living",
		"Never stop learning",
	}
	inputQuery := ""

	expected := make([][]float64, 0)

	ti := go_tfidf.New()
	err := ti.AddDocuments(inputDocuments)
	ti.CalculateDocumentsIdf()
	actual, err := ti.CalculateQueryTermsTfIdfForEachDocument(inputQuery)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func TestCalculateQueryTermsTfIdfForEachDocument(t *testing.T) {
	inputDocuments := []string{
		"The game of life is a game of everlasting learning",
		"The unexamined life is not worth living",
		"Never stop learning",
	}
	inputQuery := "life"
	expected := [][]float64{
		[]float64{0.14054651081081646},
		[]float64{0.20078072972973776},
		[]float64{0},
	}

	ti := go_tfidf.New()
	err := ti.AddDocuments(inputDocuments)
	ti.CalculateDocumentsIdf()
	actual, err := ti.CalculateQueryTermsTfIdfForEachDocument(inputQuery)

	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestCalculateQueryTermsTfIdfWhenQueryIsInvalid(t *testing.T) {
	inputQuery := ""

	expected := make([]float64, 0)

	actual, err := go_tfidf.CalculateQueryTermsTfIdf(inputQuery, " ")

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func TestCalculateQueryTermsTfIdf(t *testing.T) {
	inputQuery := "very-interesting-query"

	expected := []float64{
		0.3333333333333333,
		0.3333333333333333,
		0.3333333333333333,
	}

	actual, err := go_tfidf.CalculateQueryTermsTfIdf(inputQuery, "-")

	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

// The main reference for implementing this lib was https://janav.wordpress.com/2013/10/27/tf-idf-and-cosine-similarity/
// Thus, the tests that are going to be used are based on the tutorial example
func TestTfIdfWithCosineSimilarity(t *testing.T) {
	inputDocuments := []string{
		"The game of life is a game of everlasting learning",
		"The unexamined life is not worth living",
		"Never stop learning",
	}
	inputQuery := "life learning"

	expectedNormalizedTf := []map[string]float64{
		map[string]float64{
			"the":         0.1,
			"game":        0.2,
			"of":          0.2,
			"life":        0.1,
			"is":          0.1,
			"a":           0.1,
			"everlasting": 0.1,
			"learning":    0.1,
		},
		map[string]float64{
			"the":        0.14285714285714285,
			"unexamined": 0.14285714285714285,
			"life":       0.14285714285714285,
			"is":         0.14285714285714285,
			"not":        0.14285714285714285,
			"worth":      0.14285714285714285,
			"living":     0.14285714285714285,
		},
		map[string]float64{
			"never":    0.3333333333333333,
			"stop":     0.3333333333333333,
			"learning": 0.3333333333333333,
		},
	}
	expectedDocumentTerms := []string{
		"the",
		"game",
		"of",
		"life",
		"is",
		"a",
		"everlasting",
		"learning",
		"unexamined",
		"not",
		"worth",
		"living",
		"never",
		"stop",
	}
	expectedDocumentIdf := map[string]float64{
		"the":         1.4054651081081644,
		"game":        2.0986122886681096,
		"of":          2.0986122886681096,
		"life":        1.4054651081081644,
		"is":          1.4054651081081644,
		"a":           2.0986122886681096,
		"everlasting": 2.0986122886681096,
		"learning":    1.4054651081081644,
		"unexamined":  2.0986122886681096,
		"not":         2.0986122886681096,
		"worth":       2.0986122886681096,
		"living":      2.0986122886681096,
		"never":       2.0986122886681096,
		"stop":        2.0986122886681096,
	}
	expectedSimilarities := []float64{
		1.0,
		0.7071067811865475,
		0.7071067811865475,
	}

	ti := go_tfidf.New()
	err := ti.AddDocuments(inputDocuments)

	assert.Equal(t, expectedNormalizedTf, ti.GetDocumentsNormTermFrequency())
	assert.Equal(t, expectedDocumentTerms, ti.GetDocumentsTerms())
	assert.Nil(t, err)

	ti.CalculateDocumentsIdf()

	assert.Equal(t, expectedDocumentIdf, ti.GetDocumentsInverseFrequency())

	queryTfIdfDocuments, err := ti.CalculateQueryTermsTfIdfForEachDocument(inputQuery)
	assert.Nil(t, err)

	queryTfIdf, err := go_tfidf.CalculateQueryTermsTfIdf(inputQuery, ti.GetDocumentSeparator())
	assert.Nil(t, err)

	similarities, err := similarity.CalculateSimilarities(queryTfIdf, queryTfIdfDocuments, "Cosine")

	assert.Equal(t, expectedSimilarities, similarities)
	assert.Nil(t, err)
}
