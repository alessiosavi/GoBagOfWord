package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"math"
	"strings"

	"github.com/alessiosavi/GoBagOfWord/datastructure"
	fileutils "github.com/alessiosavi/GoGPUtils/files"
)

// dictionary will contains the words related to all document
var dictionary map[string]int = make(map[string]int)

// replacer will take in charge to remove punctation
var replacer *strings.Replacer

// initPunctation is delegated to initialize the list of punctation to remove
func initPunctation(filename string) *strings.Replacer {
	if !fileutils.FileExists(filename) {
		log.Println("File " + filename + " does not exists")
		return nil
	}

	if !fileutils.IsFile(filename) {
		log.Println("File " + filename + " is not a file :/")
		return nil
	}

	byteData, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Println("Unable to read file ...")
	}

	data := string(byteData)
	toRemove := strings.Fields(data)
	var punctation []string = make([]string, len(toRemove)*2)

	log.Println("NOTE: These symbol will be deleted -> ", toRemove)

	for i := range toRemove {
		punctation = append(punctation, toRemove[i])
		punctation = append(punctation, " ")
	}

	replacer = strings.NewReplacer(punctation...)

	return replacer
}

// initStopWords is delegated to initialize the map of stopwords
func initStopWords(filename string) {
	datastructure.StopWords = make(map[string]struct{})

	if !fileutils.FileExists(filename) {
		log.Println("File " + filename + " does not exists")
		return
	}

	if !fileutils.IsFile(filename) {
		log.Println("File " + filename + " is not a file :/")
		return
	}

	byteData, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Println("Unable to read file ...")
	}

	data := string(byteData)

	for _, field := range strings.Fields(data) {
		datastructure.StopWords[field] = struct{}{}
	}
}

//CalculateIDF is delegated to calculate the Inverse Document Frequency for each word
func CalculateIDF(docs []datastructure.DocumentData) []datastructure.DocumentData {
	var (
		// Number of total documents
		nDocument float64 = float64(len(docs))
		// Number of doc that contains the word
		wordPresent float64
	)

	// i is the index of the i-th document
	for i := range docs {
		// key is the word that we are going to analyze
		for key := range docs[i].Bow {
			wordPresent = 0
			// Check how many docs contains the word
			for j := range docs {
				if /*n*/ _, ok := docs[j].Bow[key]; ok {
					wordPresent++ // += n.Count
				}
			}
			// log.Println("Word ["+key+"] is present [", wordPresent, "] among ", nDocument, " document")
			idf := math.Log2(nDocument / wordPresent)
			_map := docs[i].Bow[key]
			_map.IDF = idf
			// (count_of_term_t_in_d) * ((log ((NUMBER_OF_DOCUMENTS + 1) / (Number_of_documents_where_t_appears +1 )) + 1)
			_map.TFIDF = _map.TF * _map.IDF
			docs[i].Bow[key] = _map
		}
		log.Println("TFIDF analyzed for [" + docs[i].DocumentName + "]")
	}
	return docs
}

//const filepath string = "/opt/DEVOPS/WORKSPACE/Golang/GoGPUtils/testdata/files/dante.txt"
const dirfolder string = "dataset"

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	var docBow []datastructure.DocumentData

	// initializing stopwords and punctation
	initStopWords("data/stopwords.txt")

	replacer = initPunctation("data/punctation.txt")

	if !fileutils.FileExists(dirfolder) {
		log.Fatal("Directory " + dirfolder + " does not exists")
	}

	// Load the document that have to be analyzed
	fileList := LoadDocumentPath(dirfolder)
	if !(len(fileList) > 0) {
		log.Fatal("Unable to find document in path [" + dirfolder + "]")
	}
	docBow = make([]datastructure.DocumentData, len(fileList))
	for i, file := range fileList {
		log.Println("Analyzing [" + file + "]")
		docBow[i].DocumentName = file
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Println("Unable to read the data for file [" + file + "]")
		} else {
			if !(len(content) > 0) {
				log.Println("File [" + file + "] is empty!")
			} else {
				// Standardizing text and calculate BoW for the document corpus
				docBow[i].Bow = StandardizeText(content, true)
			}
		}
	}

	StandardizeDict(docBow)
	docBow = CalculateIDF(docBow)
	for i := range docBow {
		log.Println(docBow[i])
	}
}

// StandardizeDict is delegated to standardize the terms that are present in all document
func StandardizeDict(docs []datastructure.DocumentData) []datastructure.DocumentData {
	var ok bool
	for i := range docs {
		for key := range dictionary {
			if _, ok = docs[i].Bow[key]; !ok {
				docs[i].Bow[key] = datastructure.BoW{}
			}
		}
	}

	return docs
}

// StandardizeText is delegated to generate the BoW for the given data
func StandardizeText(data []byte, toLower bool) map[string]datastructure.BoW {
	var (
		// This will contains the text standardized
		text string
		// Text splitted by whitespace
		words []string
		// Save the frequencies related to the word
		bow map[string]float64 = make(map[string]float64)
		// Struct for save the BoW
		bowList map[string]datastructure.BoW
	)

	if toLower {
		data = bytes.ToLower(data)
	}
	// Converting []byte in string
	text = string(data)
	text = replacer.Replace(text)
	// Tokenize text by whitespace
	words = strings.Fields(text)
	total := float64(len(words))

	// Saving the frequencies for each word
	for _, word := range words {
		// Ignore unwanted character/string
		if _, ok := datastructure.StopWords[word]; !ok {
			bow[word]++
			// Saving the words in the global dictionary
			dictionary[word] = 0
		}
	}

	// Initialize the BoW struct for save the data
	bowList = make(map[string]datastructure.BoW, len(bow))

	// Save the RAW data into the struct and calculate the TF
	for key, value := range bow {
		// if value > 0 {
		// Count -> Number of times that the terms appear
		// TF -> frequencies of the term
		bowList[key] = datastructure.BoW{Count: value, TF: value / total}
		// }
	}
	return bowList
}

// LoadDocumentPath is delegated to return the lits of file compliant with the BoW tool
func LoadDocumentPath(dirpath string) []string {
	filesList := fileutils.ListFile(dirpath)
	if filesList == nil {
		log.Fatal("Unable to load file for directory -> " + dirpath)
	}

	var files []string

	files = append(files, filesList...)

	// for _, item := range filesList {
	// 	files = append(files, item)
	// fileType, err := fileutils.GetFileContentType(item)
	// if err != nil {
	// 	log.Println("Error for file [" + item + "]")
	// } else {
	// 	if strings.HasPrefix(fileType, "text/plain") {
	// 		files = append(files, item)
	// 	} else {
	// 		log.Println("File type for file [" + item + "] -> " + fileType)
	// 	}
	// }
	// }

	return files
}
