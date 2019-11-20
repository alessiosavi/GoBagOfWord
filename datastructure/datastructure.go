package datastructure

// StopWords will contain the words that have to be ignored
var StopWords map[string]struct{}

// DocumentData is delegated to save the BoW for the document
type DocumentData struct {
	// Name of the document that we are saving the data
	DocumentName string
	// BoW of the document, the key is the word
	Bow map[string]BoW
	// TFIDF_VECTOR is delegated to save the TFIDF rappresentation of the document
	TFIDF_VECTOR []float64
}

// BoW contains the word-count for each word
type BoW struct {
	// N. of times that the word appears in the document
	Count float64
	// Frequencies of the word present into the document
	TF float64
	// Inverse document frequency related to the word present into the document
	IDF float64
	// TF * IDF
	TFIDF float64
}
