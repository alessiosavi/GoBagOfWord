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
			// Check how many docs contains the word
			for j := range docs {
				if /*n*/ _, ok := docs[j].Bow[key]; ok {
					wordPresent++ // += n.Count
				}
			}
			// log.Println("Word ["+key+"] is present [", wordPresent, "] among ", nDocument, " document")
			idf := math.Log2(nDocument + 1/wordPresent + 1) // Avoid zero division
			_map := docs[i].Bow[key]
			_map.IDF = idf
			// (count_of_term_t_in_d) * ((log ((NUMBER_OF_DOCUMENTS + 1) / (Number_of_documents_where_t_appears +1 )) + 1)
			_map.TFIDF = _map.TF * _map.IDF
			wordPresent = 0
			docs[i].Bow[key] = _map
		}
		log.Println("TFIDF analyzed for [" + docs[i].DocumentName + "]")
	}
	// for i := range docs {
	// 	log.Println("-----")
	// 	log.Println("Docs [" + docs[i].DocumentName + "]")
	// 	log.Println(docs[i])
	// 	log.Println("-----")
	// }
}

//const filepath string = "/opt/DEVOPS/WORKSPACE/Golang/GoGPUtils/testdata/files/dante.txt"
const dirfolder string = "/tmp/test"

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
	if !(len(fileList) > 0) {
		log.Fatal("Unable to find coument in path [" + dirfolder + "]")
	}
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
		//stopWords := []string{}
		docBow[i].Bow = StandardizeText(content, true, unwanted, stopWords)
	}

	// log.Println("Loaded: ", len(docBow))
	CalculateIDF(docBow)
	var wrong int
	for i := 0; i < len(docBow); i++ {
		for j := 0; j < len(docBow); j++ {
			v1, v2, penalization := FindSimilar(docBow[i], docBow[j])
			simil := CosineSimilarity(v1, v2) - penalization
			if simil > 0.5 {
				//log.Println("Document ["+docBow[i].DocumentName+"] is similar to ["+docBow[j].DocumentName+"] for a FACTOR: [", simil, "] with penalization: [", penalization, "]")
				if strings.Contains(docBow[i].DocumentName, "sci") && strings.Contains(docBow[j].DocumentName, "soc") {
					wrong++
				}
			} else {
				if strings.Contains(docBow[i].DocumentName, "sci") && strings.Contains(docBow[j].DocumentName, "sci") {
					wrong++
				} else if strings.Contains(docBow[i].DocumentName, "soc") && strings.Contains(docBow[j].DocumentName, "soc") {
					wrong++
				}
			}
		}
	}
	log.Println("Wrongs -> ", wrong)
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
		// TermFrequency frequencies of the term
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

// CosineSimilarity is delegated to calculate the Cosine Similarity for the given array
func CosineSimilarity(a, b []float64) float64 {

	if len(a) == 0 || len(b) == 0 {
		log.Fatal("CosineSimilarity | Nil input data")
	}

	if len(a) != len(b) {
		log.Fatal("CosineSimilarity | Input vectors have different size")
	}

	// Calculate numerator
	var numerator float64
	for i := range a {
		numerator += a[i] * b[i]
	}
	// Caluclate first term of denominator
	var den1 float64
	for i := range a {
		den1 += math.Pow(a[i], 2)
	}
	den1 = math.Sqrt(den1)
	// Caluclate second term of denominator
	var den2 float64
	for i := range b {
		den2 += math.Pow(b[i], 2)
	}

	den2 = math.Sqrt(den2)
	result := numerator / (den1 * den2)
	return result
}

func FindSimilar(doc1, doc2 DocumentData) ([]float64, []float64, float64) {
	var list1, list2 []float64
	var penalization float64
	for key := range doc1.Bow {
		if _, ok := doc2.Bow[key]; ok {
			//		log.Printf("Key shared! -> " + key)
			list1 = append(list1, doc1.Bow[key].Count)
			list2 = append(list2, doc2.Bow[key].Count)
		} else {
			//		log.Println("Key [" + key + "] is not shared!")
			penalization += 1 / float64(len(doc1.Bow)+len(doc2.Bow))
			//penalization++
		}
	}
	//penalization = penalization / float64(len(doc1.Bow)+len(doc2.Bow))
	// log.Println("List1-> ", list1)
	// log.Println("List2-> ", list2)

	return list1, list2, penalization
}
