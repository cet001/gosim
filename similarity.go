package gosim

// NOTES:
// - http://blog.christianperone.com/2011/09/machine-learning-text-feature-extraction-tf-idf-part-i
//

import (
	"log"
	"math"
	"os"
	"sort"
	"time"
)

var logger = log.New(os.Stderr, "[gosim] ", (log.Ldate | log.Ltime | log.Lshortfile))

type SparseHashVector map[int]float64

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

// Represents a vectorized document within a corpus.
type Document struct {
	// Unique document ID within a given corpus.
	Id int

	// Term frequencies for each unique term in this document.
	tf SparseVector

	// TF-IDF score of each distinct term x in this document.
	tfidf SparseVector
}

// TF-IDF model.
type TFIDF struct {
	// The documents within this corpus.
	docs []Document

	// idf[x] -> the inverse document frequency of term x.
	idf SparseHashVector

	// Whenever new documents are added to this corpus, the global stats need to
	// be recalculated (via Recalculate()).  This flag keeps track of this state.
	needsRecalc bool
}

func NewTFIDF() *TFIDF {
	return &TFIDF{
		docs:        make([]Document, 0, 200000),
		needsRecalc: true,
	}
}

func (me *TFIDF) AddDoc(docId int, doc SparseVector) {
	me.docs = append(me.docs, Document{
		Id: docId,
		tf: doc,
	})
	me.needsRecalc = true
}

// Trains the model. Returns a list of the distinct terms and their
// corresponding document frequency (sorted by increasing frequency).
func (me *TFIDF) Train() []Term {
	logger.Printf("Calculating document frequencies")
	startTime := time.Now()
	df := calcDocFrequencies(me.docs)
	logger.Printf("Document frequency calculation took %v.", time.Since(startTime))

	logger.Printf("Removing unimportant terms from corpus")
	startTime = time.Now()
	unimportantTerms := removeUnimportantTerms(df, len(me.docs))
	filterDocVectors(me.docs, df)
	logger.Printf("%v unimportant terms removed in %v.", len(unimportantTerms), time.Since(startTime))

	logger.Printf("Calculating IDF values for %v terms.", len(df))
	startTime = time.Now()
	me.idf = make(SparseHashVector, len(df))
	totalDocs := float64(len(me.docs))
	for termId, df := range df {
		me.idf[termId] = 1.0 + math.Log(totalDocs/float64(df))
	}
	logger.Printf("IDF calculation took %v.", time.Since(startTime))

	logger.Printf("Calculating TF-IDF values")
	startTime = time.Now()
	for i := 0; i < len(me.docs); i++ {
		doc := &me.docs[i]
		doc.tfidf = calcTFIDF(doc.tf, me.idf)
	}
	logger.Printf("TF-IDF calculation took %v.", time.Since(startTime))

	me.needsRecalc = false
	return unimportantTerms
}

// Calculates a similarity score indicating how similar douments doc1 and doc2 are.
// Returns a score in the range [0.0..1.0], where 1.0 means the documents are
// identical.
func (me *TFIDF) CalcSimilarity(doc1, doc2 SparseVector) float64 {
	me.validateState()

	// Calculate cosine similarity
	doc1_tfidf := calcTFIDF(doc1, me.idf)
	doc2_tfidf := calcTFIDF(doc2, me.idf)
	return Dot(doc1_tfidf, doc2_tfidf) / (Norm(doc1_tfidf) * Norm(doc2_tfidf))
}

// Ranks the documents in the corpus in terms of how similar they are to the
// specified query.
func (me *TFIDF) SimilarDocsForText(query SparseVector) []ScoredItem {
	me.validateState()

	queryTFIDF := calcTFIDF(query, me.idf)
	normQueryTFIDF := Norm(queryTFIDF)

	rankedDocs := []ScoredItem{}
	for i := 0; i < len(me.docs); i++ {
		doc := &me.docs[i]

		if len(doc.tfidf) > 0 {
			score := Dot(queryTFIDF, doc.tfidf) / (normQueryTFIDF * Norm(doc.tfidf))
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

// Calculates the document frequency for each distinct term contained in the
// specified list of documents.
func calcDocFrequencies(docs []Document) map[int]int {
	df := make(map[int]int, 1000000)
	for i := 0; i < len(docs); i++ {
		doc := &docs[i]
		for j := 0; j < len(doc.tf); j++ {
			term := &doc.tf[j]
			df[term.Id] += 1
		}
	}
	return df
}

// termFreqs = term frequencies
// idfs = vector of inverse document frequencies for each term in the corpus.
func calcTFIDF(termFreqs SparseVector, idfs SparseHashVector) SparseVector {
	tfidf := make([]Term, len(termFreqs))
	for i := 0; i < len(termFreqs); i++ {
		term := &termFreqs[i]
		tfidf[i] = Term{Id: term.Id, Value: (term.Value * idfs[term.Id])}
	}

	return tfidf
}

func removeUnimportantTerms(docFreqs map[int]int, numDocs int) []Term {
	removedTerms := make([]Term, 0, 100000)

	for termId, docFreq := range docFreqs {
		isUnimportantTerm := (docFreq <= 3 || ((float64(docFreq) / float64(numDocs)) > 0.20))
		if isUnimportantTerm {
			delete(docFreqs, termId)
			removedTerms = append(removedTerms, Term{Id: termId, Value: float64(docFreq)})
		}
	}

	return removedTerms
}

// For each document, keeps *only* the term Ids within that document's term frequency
// vector (doc.tf) that are present in the provided termLookup map.
// The keys in the termLookup map represent the term Ids, and the values are actually
// not read by this function.
func filterDocVectors(docs []Document, termLookup map[int]int) {
	for i := 0; i < len(docs); i++ {
		doc := &docs[i]

		filteredVec := make(SparseVector, 0, len(doc.tf))
		for _, term := range doc.tf {
			_, found := termLookup[term.Id]
			if found {
				filteredVec = append(filteredVec, term)
			}
		}
		doc.tf = filteredVec
	}
}
