// Package tfidf provides an implementation of the TF-IDF statistical model.
//
// See https://en.wikipedia.org/wiki/Tf-idf
package tfidf

// NOTES:
// - http://blog.christianperone.com/2011/09/machine-learning-text-feature-extraction-tf-idf-part-i
//

import (
	"encoding/gob"
	"github.com/cet001/gosim/math"
	"log"
	gomath "math"
	"os"
	"sort"
	"time"
)

var logger = log.New(os.Stderr, "[gosim] ", (log.Ldate | log.Ltime | log.Lshortfile))

// Represents a vectorized document within a corpus.
type Document struct {
	// Unique document ID within a given corpus.
	Id int

	// Term frequencies for each unique term in this document.
	TF math.SparseVector

	// TF-IDF score of each distinct term x in this document.
	TFIDF math.SparseVector
}

// Statistics that were gathered during the training phase (see Train()).
type Stats struct {
	// The number of documents in the corpus
	DocumentCount int

	// The number of distinct terms in the corpus
	TermCount int

	// The stop words that were identified by this algorithm.
	StopWords []math.Term
}

// Anything that can be represented as a unique Id and associated score.
type ScoredItem struct {
	Id    int
	Score float64
}

// Sorts ScoredItem objects in descending order by score.
type byScore []ScoredItem

func (a byScore) Len() int           { return len(a) }
func (a byScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byScore) Less(i, j int) bool { return a[i].Score > a[j].Score }

// Sparse vector represented by a mapping of term IDs to corresponding term values.
// This is essentially the hashmap version of the SparseVector type.
type sparseHashVector map[int]float64

// TF-IDF model.
type TFIDF struct {
	// A term will be considered a stopword if it is present in more than the
	// percentage of documents in the corpus specified by this field.  Valid
	// range is [0..1], where 0 = 0% and 1 = 100%.
	StopWordThreshold float64

	// The documents within this corpus.
	docs []Document

	// idf[t] -> the inverse document frequency of term t.
	idf sparseHashVector

	// Whenever new documents are added to this corpus, the global stats need to
	// be recalculated (via Recalculate()).  This flag keeps track of this state.
	needsRecalc bool
}

func NewTFIDF() *TFIDF {
	return &TFIDF{
		StopWordThreshold: 0.20,
		docs:              make([]Document, 0, 200000),
		needsRecalc:       true,
	}
}

// Saves this model to the specified file.
func (me *TFIDF) Save(filePath string) error {
	file, err := os.Create(filePath)
	defer file.Close()

	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(me.StopWordThreshold)
		encoder.Encode(len(me.docs))
		for i := 0; i < len(me.docs); i++ {
			encoder.Encode(&me.docs[i])
		}
	}

	return err
}

// Loads a TFIDF model from a saved image on file.
func LoadTFIDF(filePath string) (*TFIDF, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)

	var stopWordThreshold float64
	if err := decoder.Decode(&stopWordThreshold); err != nil {
		return nil, err
	}

	var docCount int
	if err := decoder.Decode(&docCount); err != nil {
		return nil, err
	}

	docs := make([]Document, 0, docCount)
	for i := 0; i < docCount; i++ {
		var doc Document
		if err := decoder.Decode(&doc); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	return &TFIDF{
		StopWordThreshold: stopWordThreshold,
		docs:              docs,
		needsRecalc:       true,
	}, nil
}

func (me *TFIDF) AddDoc(docId int, doc math.SparseVector) {
	me.docs = append(me.docs, Document{
		Id: docId,
		TF: doc,
	})
	me.needsRecalc = true
}

// Trains the model. Returns a list of the distinct terms and their
// corresponding document frequency (sorted by increasing frequency).
func (me *TFIDF) Train() Stats {
	logger.Printf("Calculating document frequencies")
	startTime := time.Now()
	df := calcDocFrequencies(me.docs)
	logger.Printf("Document frequency calculation took %v.", time.Since(startTime))

	logger.Printf("Removing stop words from document frequency map")
	startTime = time.Now()
	stopWords := removeStopWords(df, len(me.docs), me.StopWordThreshold)
	logger.Printf("%v stop words removed in %v.", len(stopWords), time.Since(startTime))

	logger.Printf("Filtering document vectors based on reduced document frequency map")
	startTime = time.Now()
	filterDocVectors(me.docs, df)
	logger.Printf("Document vector filtering took %v.", time.Since(startTime))

	logger.Printf("Calculating IDF values for %v terms.", len(df))
	startTime = time.Now()
	me.idf = make(sparseHashVector, len(df))
	totalDocs := float64(len(me.docs))
	for termId, df := range df {
		me.idf[termId] = 1.0 + gomath.Log(totalDocs/float64(df))
	}
	logger.Printf("IDF calculation took %v.", time.Since(startTime))

	logger.Printf("Calculating TF-IDF values")
	startTime = time.Now()
	for i := 0; i < len(me.docs); i++ {
		doc := &me.docs[i]
		doc.TFIDF = calcTFIDF(doc.TF, me.idf)
	}
	logger.Printf("TF-IDF calculation took %v.", time.Since(startTime))

	me.needsRecalc = false
	return Stats{
		DocumentCount: len(me.docs),
		TermCount:     len(df),
		StopWords:     stopWords,
	}
}

// Calculates a similarity score indicating how similar documents doc1 and doc2
// to each other.
//
// Returns a score in the range [0.0..1.0], where 1.0 means the documents are
// identical.
func (me *TFIDF) CalcSimilarity(doc1, doc2 math.SparseVector) float64 {
	me.validateState()

	// Calculate cosine similarity
	doc1_tfidf := calcTFIDF(doc1, me.idf)
	doc2_tfidf := calcTFIDF(doc2, me.idf)
	score := math.Dot(doc1_tfidf, doc2_tfidf) / (math.Norm(doc1_tfidf) * math.Norm(doc2_tfidf))
	return gomath.Min(1.0, score)
}

// Ranks the documents in the corpus in terms of how similar they are to the
// specified query.
func (me *TFIDF) SimilarDocsForText(query math.SparseVector) []ScoredItem {
	me.validateState()

	queryTFIDF := calcTFIDF(query, me.idf)
	normQueryTFIDF := math.Norm(queryTFIDF)

	rankedDocs := []ScoredItem{}
	for i := 0; i < len(me.docs); i++ {
		doc := &me.docs[i]

		if len(doc.TFIDF) > 0 {
			score := math.Dot(queryTFIDF, doc.TFIDF) / (normQueryTFIDF * math.Norm(doc.TFIDF))
			score = gomath.Min(1.0, score)
			rankedDocs = append(rankedDocs, ScoredItem{Id: doc.Id, Score: score})
		}
	}

	sort.Sort(byScore(rankedDocs))
	return rankedDocs
}

// Call this method to ensure the corpus is in state that it can be queried.
func (me *TFIDF) validateState() {
	if me.needsRecalc {
		panic("Corpus stats need to be recalculated.  Call Recalculate().")
	}
}

// Calculates the document frequency (df) for each distinct term within the
// specified corpus.  Returns a map df[t], where t is a distinct term ID,
// and df[t] returns the number of documents in the corpus that contain at least
// one mention of t.
func calcDocFrequencies(corpus []Document) map[int]int {
	df := make(map[int]int, 1000000)

	for i := 0; i < len(corpus); i++ {
		doc := &corpus[i]
		for j := 0; j < len(doc.TF); j++ {
			term := &doc.TF[j]
			df[term.Id] += 1
		}
	}

	return df
}

// termFreqs = term frequencies
// idfs = vector of inverse document frequencies for each term in the corpus.
func calcTFIDF(termFreqs math.SparseVector, idfs sparseHashVector) math.SparseVector {
	tfidf := make([]math.Term, len(termFreqs))
	for i := 0; i < len(termFreqs); i++ {
		term := &termFreqs[i]
		tfidf[i] = math.Term{
			Id:    term.Id,
			Value: (term.Value * idfs[term.Id]),
		}
	}

	return tfidf
}

// Identifies stopwords within the specified docFreqs map and then removes them.
// A stopword is defined as a word that is present in more than 20% of the
// documents in the corpus.
func removeStopWords(docFreqs map[int]int, numDocs int, threshold float64) []math.Term {
	stopWords := make([]math.Term, 0, 100000)
	for termId, docFreq := range docFreqs {
		isStopWord := (float64(docFreq) / float64(numDocs)) > threshold
		if isStopWord {
			delete(docFreqs, termId)
			stopWords = append(stopWords, math.Term{Id: termId, Value: float64(docFreq)})
		}
	}

	return stopWords
}

// Identifies and removes "unimportant" terms within the given document
// frequency map.
func removeUnimportantTerms(docFreqs map[int]int) []math.Term {
	removedTerms := make([]math.Term, 0, 100000)

	for termId, docFreq := range docFreqs {
		isRareInCorpus := (docFreq <= 1)
		if isRareInCorpus {
			delete(docFreqs, termId)
			removedTerms = append(removedTerms, math.Term{Id: termId, Value: float64(docFreq)})
		}
	}

	return removedTerms
}

// For each document in []docs, keeps *only* the term Ids within that document's
// term frequency vector (doc.tf) that are present in the provided termLookup
// map.  The keys in the termLookup map represent the term Ids, and the values
// are not used by this function.
func filterDocVectors(docs []Document, termLookup map[int]int) {
	for i := 0; i < len(docs); i++ {
		doc := &docs[i]

		filteredVec := make(math.SparseVector, 0, len(doc.TF))
		for _, term := range doc.TF {
			_, found := termLookup[term.Id]
			if found {
				filteredVec = append(filteredVec, term)
			}
		}
		doc.TF = filteredVec
	}
}
