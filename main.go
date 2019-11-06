package main

import (
	"bytes"
	"io/ioutil"
	"log"
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
	// N. of times that the word appears
	Count float64
	// Frequencies in relation to the document
	TermFrequency float64
}

//CalculateIDF is delegated to calculate the Inverse Document Frequency for each word
func CalculateIDF(docs DocumentData) {

}

//const filepath string = "/opt/DEVOPS/WORKSPACE/Golang/GoGPUtils/testdata/files/dante.txt"
const dirfolder string = "dataset/"

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

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

		docBow[i].DocumentName = file
		content, err := ioutil.ReadFile(file)

		if err != nil {
			log.Fatal("Unable to read the data for ->" + file)
		}

		unwanted = []string{",", ":", ";", ".", "‘", "”", "“", "»", "«", "<<", ">>", "?", "!"}
		docBow[i].Bow = StandardizeText(content, true, unwanted)
	}

	log.Println("Loaded: ", len(docBow))
}

// StandardizeText is delegated to generate the BoW for the given data
func StandardizeText(data []byte, toLower bool, toRemove []string) map[string]BoW {
	var (
		// This will contains the text
		text string
		// Text splitted by whitespace
		words []string
		// Save the frequencies related to the word
		bow map[string]float64 = make(map[string]float64)
		// Struct for save the BoW and TF
		bowList map[string]BoW = make(map[string]BoW)
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
			unwanted = append(unwanted, "")
		}

		replacer := strings.NewReplacer(unwanted...)
		text = replacer.Replace(text)
	}
	// Split the text
	words = strings.Fields(text)

	total := float64(len(words))
	// Saving the frequencies for each word
	for _, word := range words {
		bow[word]++
	}

	// Initialize the BoW struct for save the data
	bowList = make(map[string]BoW, len(bow))

	// Save the RAW data into the struct and calculate the TF
	for key, value := range bow {
		bowList[key] = BoW{Count: value, TermFrequency: value / total}
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
		fileType, err := fileutils.GetFileContentType(item)
		if err != nil {
			log.Println("Error for file [" + item + "]")
		} else {
			if strings.HasPrefix(fileType, "text/plain") {
				files = append(files, item)
			} else {
				log.Println("File type for file [" + item + "] -> " + fileType)
			}
		}
	}
	return files
}
