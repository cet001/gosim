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

type IntVector map[int]int

type FloatVector map[int]float64

// Anything with a unique ID and a score.
type Term struct {
	Id    int
	Score float64
}

// Sorts terms in descending order by score.
type byScore []Term

func (a byScore) Len() int           { return len(a) }
func (a byScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byScore) Less(i, j int) bool { return a[i].Score > a[j].Score }

// Sorts terms by increasing Id.
type byTermId []Term

func (a byTermId) Len() int           { return len(a) }
func (a byTermId) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byTermId) Less(i, j int) bool { return a[i].Id < a[j].Id }

// Represents a vectorized document within a corpus.
type Document struct {
	// Unique document ID within a given corpus.
	Id int

	// Term frequencies for each unique term in this document.
	tf []Term

	// tfidf[x] -> the TF-IDF score of term x within this document.
	tfidf []Term
}

type TFIDF struct {
	// The documents within this corpus.
	docs []Document

	// df[x] -> the number of documents in which term x was mentioned.
	df IntVector

	// idf[x] -> the inverse document frequency of term x.
	idf FloatVector

	// Whenever new documents are added to this corpus, the global stats need to
	// be recalculated (via Recalculate()).  This flag keeps track of this state.
	needsRecalc bool
}

func NewTFIDF() *TFIDF {
	return &TFIDF{
		docs:        make([]Document, 0, 200000),
		df:          make(IntVector, 200000),
		needsRecalc: true,
	}
}

func (me *TFIDF) AddDoc(docId int, termFrequencies []Term) {
	me.docs = append(me.docs, Document{
		Id: docId,
		tf: termFrequencies,
	})
	me.needsRecalc = true
}

func (me *TFIDF) Calculate() {
	// const termFreqThreshold = 2
	// logger.Printf("Trimming vocabulary to those terms with fewer than %v doc mentions.", termFreqThreshold)
	// startTime := time.Now()
	// totalTermCount := len(me.df)
	// unimportantTermCount := 0
	// for term, termId := range me.vocab {
	// 	if me.df[termId] < termFreqThreshold {
	// 		unimportantTermCount++
	// 		delete(me.vocab, term)
	// 		delete(me.df, termId)
	// 	}
	// }
	// logger.Printf("%v of %v terms occurred in fewer than %v documents.", unimportantTermCount, totalTermCount, termFreqThreshold)
	// for _, doc := range me.id2doc {
	// 	for termId, _ := range doc.tf {
	// 		_, isVocabTerm := me.df[termId]
	// 		if !isVocabTerm {
	// 			delete(doc.tf, termId)
	// 		}
	// 	}
	// }
	// logger.Printf("Vocabulary trimming took %v.", time.Since(startTime))

	logger.Printf("Calculating document frequencies")
	startTime := time.Now()
	df := make(map[int]int, 1000000)
	for i := 0; i < len(me.docs); i++ {
		doc := &me.docs[i]
		for j := 0; j < len(doc.tf); j++ {
			term := &doc.tf[j]
			df[term.Id] += 1
		}
	}
	logger.Printf("Document frequency calculation took %v.", time.Since(startTime))

	logger.Printf("Calculating IDF values for %v terms.", len(df))
	startTime = time.Now()
	me.idf = make(FloatVector, len(df))
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
}

func (me *TFIDF) CalcSimilarity(doc1, doc2 []Term) float64 {
	me.validateState()

	doc1_tfidf := calcTFIDF(doc1, me.idf)
	doc2_tfidf := calcTFIDF(doc2, me.idf)
	return dot(doc1_tfidf, doc2_tfidf) / (norm(doc1_tfidf) * norm(doc2_tfidf))
}

func (me *TFIDF) SimilarDocsForText(query []Term) []Term {
	me.validateState()

	queryTFIDF := calcTFIDF(query, me.idf)
	normQueryTFIDF := norm(queryTFIDF)
	logger.Printf(">>>> SimilarDocsForText(): queryTFIDF vector size is %v", len(queryTFIDF))

	rankedDocs := []Term{}
	for i := 0; i < len(me.docs); i++ {
		doc := &me.docs[i]
		score := dot(queryTFIDF, doc.tfidf) / (normQueryTFIDF * norm(doc.tfidf))
		rankedDocs = append(rankedDocs, Term{Id: doc.Id, Score: score})
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

// doc = Document represented as a vector of (term, freq) tuples
// idfs = vector of pre-calculated inverse document frequencies for each term in the corpus.
func calcTFIDF(doc []Term, idfs FloatVector) []Term {
	tfidf := make([]Term, len(doc))
	for i := 0; i < len(doc); i++ {
		term := &doc[i]
		tfidf[i] = Term{Id: term.Id, Score: (term.Score * idfs[term.Id])}
	}

	return tfidf
}

// Calculates the Euclidean norm (a.k.a. L2-norm) of the specified vector.
func norm(vec []Term) float64 {
	sumOfSquares := 0.0
	for i := 0; i < len(vec); i++ {
		term := &vec[i]
		sumOfSquares += (term.Score * term.Score)
	}

	return math.Sqrt(sumOfSquares)
}

// Calculates the dot product of vectors v1 and v2.
func dot(v1, v2 []Term) float64 {
	var dp float64 = 0.0
	lenV1, lenV2 := len(v1), len(v2)
	idx1, idx2 := 0, 0

	for {
		term1, term2 := &v1[idx1], &v2[idx2]

		if term1.Id < term2.Id {
			idx1++
		} else if term2.Id < term1.Id {
			idx2++
		} else {
			dp += (term1.Score * term2.Score)
			idx1++
			idx2++
		}

		if idx1 == lenV1 || idx2 == lenV2 {
			break
		}
	}

	return dp
}
