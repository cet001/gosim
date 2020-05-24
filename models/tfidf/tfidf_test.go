package tfidf

import (
	"fmt"
	"github.com/cet001/gosim"
	"github.com/cet001/mathext/vectors"
	"github.com/stretchr/testify/assert"
	"os"
	"sort"
	"testing"
)

func ExampleTFIDF() {
	// Our corpus
	corpus := []string{
		"Life is about making an impact, not making an income.",
		"Whatever the mind of man can conceive and believe, it can achieve.",
		"Strive not to be a success, but rather to be of value.",
		"Two roads diverged in a wood, and I —- I took the one less traveled by, And that has made all the difference",
		"I attribute my success to this: I never gave or took any excuse.",
		"You miss 100 percent of the shots you don't take.",
		"I've missed more than 9000 shots in my career. I've lost almost 300 games. 26 times I've been trusted to take the game winning shot and missed. I've failed over and over and over again in my life. And that is why I succeed.",
		"The most difficult thing is the decision to act, the rest is merely tenacity.",
		"Every strike brings me closer to the next home run.",
		"Definiteness of purpose is the starting point of all achievement.",
	}

	// Initialize the TFIDF model
	model := NewTFIDF()
	dict := gosim.NewDictionary()
	tokenize := gosim.MakeDefaultTokenizer()

	// Vectorize each document and then insert it into our TFIDF model
	for docId, doc := range corpus {
		words := tokenize(doc)
		docVector := dict.VectorizeAndUpdate(words)
		model.AddDoc(docId, docVector)
	}

	stats := model.Train()

	// Resolve the StopWord terms to their string token values and then sort them.
	stopWords := []string{}
	for _, stopWord := range stats.StopWords {
		stopWords = append(stopWords, dict.Word(stopWord.Id))
	}
	sort.Strings(stopWords)

	fmt.Printf("Unique terms in corpus: %v\n", stats.TermCount)
	fmt.Printf("Stop words: %v\n", stopWords)

	// Output:
	// Unique terms in corpus: 94
	// Stop words: [and is of the to]
}

func TestSaveAndLoadTFIDF(t *testing.T) {
	doc := Document{
		Id:    123,
		TF:    vectors.SparseVector{{1, 10}, {2, 20}},
		TFIDF: vectors.SparseVector{{1, 0.5}, {2, 0.25}},
	}
	docs := []Document{doc}

	model := TFIDF{
		StopWordThreshold: 0.20,
		docs:              docs,
	}

	dataFile := "/tmp/gosim_TestSaveAndLoadModel.dat"
	defer os.Remove(dataFile)

	saveErr := model.Save(dataFile)
	assert.Nil(t, saveErr)

	reloadedModel, loadErr := LoadTFIDF(dataFile)
	assert.Nil(t, loadErr)
	assert.Equal(t, reloadedModel.StopWordThreshold, model.StopWordThreshold)
	assert.Equal(t, reloadedModel.docs, model.docs)
}

func TestLoadTFIDF_nonexistentFile(t *testing.T) {
	_, err := LoadTFIDF("/a/b/c/nonexistent-file-xxxxxxxxxx.txt")
	assert.NotNil(t, err)
}

func TestCalcTFIDF(t *testing.T) {
	tf := vectors.SparseVector{{10, 10}, {40, 40}, {50, 50}}
	idf := sparseHashVector{10: 0.1, 20: 0.2, 30: 0.3, 40: 0.4, 50: 0.5}
	assert.Equal(t,
		vectors.SparseVector{{10, (10 * 0.1)}, {40, (40 * 0.4)}, {50, (50 * 0.5)}},
		calcTFIDF(tf, idf),
	)
}

func TestTFIDF_AddDoc(t *testing.T) {
	c := NewTFIDF()

	docvec1 := []vectors.Element{
		{Id: 1, Value: 10},
		{Id: 2, Value: 20},
	}

	docvec2 := []vectors.Element{
		{Id: 1, Value: 100},
		{Id: 3, Value: 300},
	}

	// Add 2 document vectors
	c.needsRecalc = false
	c.AddDoc(10000, docvec1)
	c.AddDoc(20000, docvec2)

	assert.Equal(t, 2, len(c.docs))
	assert.True(t, c.needsRecalc)
}

func TestTFIDF_CalcSimilarity(t *testing.T) {
	corpus := []string{
		"apache helicopter military war", // docId=0
		"war helicopter apache military", // docId=1
		"apache software code developer", // docId=2
		"foo bar baz",                    // docId=3
	}

	// Initialize the TFIDF model
	model := NewTFIDF()
	model.StopWordThreshold = 1.0 // effectively turns off stopword filtering
	dict := gosim.NewDictionary()
	tokenize := gosim.MakeDefaultTokenizer()

	// Vectorize each document and then insert it into our TFIDF model
	vectorizedDocs := []vectors.SparseVector{}
	for docId, doc := range corpus {
		words := tokenize(doc)
		docVector := dict.VectorizeAndUpdate(words)
		model.AddDoc(docId, docVector)
		vectorizedDocs = append(vectorizedDocs, docVector)
	}

	model.Train()

	assert.Equal(t, 1.0, model.CalcSimilarity(vectorizedDocs[0], vectorizedDocs[1]))
	assert.Equal(t, 0.0, model.CalcSimilarity(vectorizedDocs[0], vectorizedDocs[3]))
}

// This is just a sanity check
func TestTFIDF_SimilarDocsForText(t *testing.T) {
	corpus := []string{
		"apache helicopter military war", // docId=0
		"apache software code developer", // docId=1
		"apache indian history tribe",    // docId=2
	}

	// Initialize the TFIDF model
	model := NewTFIDF()
	model.StopWordThreshold = 1.0 // effectively turns off stopword filtering
	dict := gosim.NewDictionary()
	tokenize := gosim.MakeDefaultTokenizer()

	// Vectorize each document and then insert it into our TFIDF model
	for docId, doc := range corpus {
		words := tokenize(doc)
		docVector := dict.VectorizeAndUpdate(words)
		model.AddDoc(docId, docVector)
	}

	model.Train()

	plainTextQuery := "apache http server is used by countless software projects"
	tokenizedQuery := tokenize(plainTextQuery)
	vectorizedQuery := dict.Vectorize(tokenizedQuery)
	similarDocs := model.SimilarDocsForText(vectorizedQuery)

	mostSimilarDoc := similarDocs[0]
	assert.Equal(t, 1, mostSimilarDoc.Id)
}

func TestCalcDocFrequencies(t *testing.T) {
	t1, t2, t3 := vectors.Element{10, 1}, vectors.Element{20, 1}, vectors.Element{30, 1}
	docs := []Document{
		{
			Id: 100,
			TF: vectors.SparseVector{t1, t2, t3},
		},
		{
			Id: 200,
			TF: vectors.SparseVector{t2, t3},
		},
		{
			Id: 300,
			TF: vectors.SparseVector{t3},
		},
	}

	assert.Equal(t, map[int]int{10: 1, 20: 2, 30: 3}, calcDocFrequencies(docs))
}

func TestRemoveStopWords(t *testing.T) {
	docCount := 10
	docFreqs := map[int]int{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
		5: 5,
	}
	removedTerms := removeStopWords(docFreqs, docCount, 0.20)
	assert.Equal(t, map[int]int{1: 1, 2: 2}, docFreqs)
	assert.Equal(t, 3, len(removedTerms))
}

func TestRemoveUnimportantTerms(t *testing.T) {
	docFreqs := map[int]int{
		1: 1,
		2: 222,
		3: 1,
		4: 444,
		5: 555,
	}
	removedTerms := removeUnimportantTerms(docFreqs)
	assert.Equal(t, map[int]int{2: 222, 4: 444, 5: 555}, docFreqs)
	assert.Equal(t, 2, len(removedTerms))
}

func TestFilterDocVectors(t *testing.T) {
	docs := []Document{
		{
			Id: 100,
			TF: vectors.SparseVector{{10, 1}, {20, 2}, {30, 3}},
		},
		{
			Id: 200,
			TF: vectors.SparseVector{{20, 200}, {30, 300}, {40, 400}},
		},
	}

	// Create a filter that specifies which term Ids are to be kept
	someValue := 999
	filter := map[int]int{10: someValue, 30: someValue, 50: someValue}

	filterDocVectors(docs, filter)

	assert.Equal(t,
		[]Document{
			{
				Id: 100,
				TF: vectors.SparseVector{{10, 1}, {30, 3}},
			},
			{
				Id: 200,
				TF: vectors.SparseVector{{30, 300}},
			},
		},
		docs,
	)
}
