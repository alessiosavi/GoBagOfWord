package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"math"
	"strings"

	fileutils "github.com/alessiosavi/GoGPUtils/files"
)

// DocumentData is delegated to save the BoW for the document
type DocumentData struct {
	// Name of the document that we are saving the data
	DocumentName string
	// BoW of the document, the key is the word
	Bow map[string]BoW
}

// BoW contains the word-count for each word
type BoW struct {
	// N. of times that the word appears in the document
	Count float64
	// Frequencies in relation to the document
	TF float64
	// Inverse document frequency related to all document
	IDF float64
	// TF * IDF
	TFIDF float64
}

//CalculateIDF is delegated to calculate the Inverse Document Frequency for each word
func CalculateIDF(docs []DocumentData) {
	// Number of total documents
	var nDocument float64
	// Number of doc that contains the word
	var wordPresent float64
	nDocument = float64(len(docs))
	//  y := math.Log(2.7183)

	// i is the index of the document
	for i := range docs {
		log.Println("Analzying TFIDF [" + docs[i].DocumentName + "]")
		// key is the word that we are going to analyze
		for key := range docs[i].Bow {
			// Check how many time this word is present
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

			_map.TFIDF = _map.IDF * _map.TF
			docs[i].Bow[key] = _map
			wordPresent = 0
		}
	}

	for i := range docs {
		log.Println("-----")
		log.Println("Docs [" + docs[i].DocumentName + "]")
		log.Println(docs[i])
		log.Println("-----")
	}

}

//const filepath string = "/opt/DEVOPS/WORKSPACE/Golang/GoGPUtils/testdata/files/dante.txt"
const dirfolder string = "dataset"

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	var (
		docBow   []DocumentData
		unwanted []string
	)

	if !fileutils.FileExists(dirfolder) {
		log.Fatal("Directory " + dirfolder + " does not exists")
	}

	fileList := LoadDocumentPath(dirfolder)
	docBow = make([]DocumentData, len(fileList))
	for i, file := range fileList {
		log.Println("Analyzing [" + file + "]")
		docBow[i].DocumentName = file
		content, err := ioutil.ReadFile(file)

		if err != nil {
			log.Fatal("Unable to read the data for ->" + file)
		}
		unwanted = []string{",", ":", ";", ".", "‘", "”", "“", "+", "»", "«", "<<", ">>", "?", "!"}
		stopWords := []string{`i`, " me ", " my ", " myself ", " we ", " our ", " ours ", " ourselves ", " you ", " you're ", " you've ", " you'll ", " you'd ", " your ", " yours ", " yourself ", " yourselves ", " he ", " him ", " his ", " himself ", " she ", " she's ", " her ", " hers ", " herself ", " it ", " it's ", " its ", " itself ", " they ", " them ", " their ", " theirs ", " themselves ", " what ", " which ", " who ", " whom ", " this ", " that ", "that'll", " these ", " those ", " am ", " is ", " are ", " was ", " were ", " be ", " been ", " being ", " have ", " has ", " had ", " having ", " do ", " does ", " did ", " doing ", `a`, " an ", " the ", " and ", " but ", `if`, `or`, `because`, `as`, `until`, `while`, `of`, `at`, `by`, `for`, `with`, `about`, `against`, `between`, `into`, `through`, `during`, `before`, `after`, `above`, `below`, `to`, `from`, `up`, `down`, `in`, `out`, `on`, `off`, `over`, `under`, `again`, `further`, `then`, `once`, `here`, `there`, `when`, `where`, `why`, `how`, `all`, `any`, `both`, `each`, `few`, `more`, `most`, `other`, `some`, `such`, `no`, `nor`, `not`, `only`, `own`, `same`, `so`, `than`, `too`, `very`, `s`, `t`, `can`, `will`, `just`, `don`, "don`t", `should`, "should`ve", `now`, `d`, `ll`, `m`, `o`, `re`, `ve`, `y`, `ain`, `aren`, "aren`t", `couldn`, "couldn`t", `didn`, "didn`t", `doesn`, "doesn`t", `hadn`, "hadn`t", `hasn`, "hasn`t", `haven`, "haven`t", `isn`, "isn`t", `ma`, `mightn`, "mightn`t", `mustn`, "mustn't", `needn`, "needn't", `shan`, "shan't", `shouldn`, "shouldn't", `wasn`, "wasn't", `weren`, "weren't", `won`, "won't", `wouldn`, "wouldn't"}

		docBow[i].Bow = StandardizeText(content, true, unwanted, stopWords)
	}

	// log.Println("Loaded: ", len(docBow))
	CalculateIDF(docBow)
}

// StandardizeText is delegated to generate the BoW for the given data
func StandardizeText(data []byte, toLower bool, toRemove, stopWords []string) map[string]BoW {
	var (
		// This will contains the text
		text string
		// Text splitted by whitespace
		words []string
		// Save the frequencies related to the word
		bow map[string]float64 = make(map[string]float64)
		// Struct for save the BoW and TF
		bowList map[string]BoW
	)

	if toLower {
		data = bytes.ToLower(data)
	}
	// Converting []byte in string
	text = string(data)

	// Removing unwanted character/string
	if len(toRemove) > 0 {
		var unwanted []string

		for i := range toRemove {
			unwanted = append(unwanted, toRemove[i])
			unwanted = append(unwanted, " ")
		}

		replacer := strings.NewReplacer(unwanted...)
		text = replacer.Replace(text)
	}
	// Split the text
	words = strings.Fields(text)
	total := float64(len(words))
	var isStopWord bool
	// Saving the frequencies for each word
	for _, word := range words {
		for i := range stopWords {
			if stopWords[i] == word {
				isStopWord = true
				break
			}
		}
		if !isStopWord {
			bow[word]++
		}
		isStopWord = false
	}

	// Initialize the BoW struct for save the data
	bowList = make(map[string]BoW, len(bow))

	// Save the RAW data into the struct and calculate the TF
	for key, value := range bow {
		// Count -> Number of times that the terms appear
		// TermFrequency
		bowList[key] = BoW{Count: value, TF: value / total}
		// log.Println("Bow -> ", bowList[key], " Key: ", key)
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
	for _, item := range filesList {
		files = append(files, item)
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
	}
	return files
}
