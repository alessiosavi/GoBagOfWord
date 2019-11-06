package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	fileutils "github.com/alessiosavi/GoGPUtils/files"
)

// DocumentData is delegated to save the BoW for the document
type DocumentData struct {
	// Name of the document that we are saving the data
	DocumentName string
	// BoW of the document
	Bow []BoW
}

// BoW contains the word-count for each word
type BoW struct {
	// Word that we are saving information in this struct
	Word string
	// N. of times that the word appears
	Count float64
	// Frequencies in relation to the document
	TermFrequency float64
}

const filepath string = "/opt/DEVOPS/WORKSPACE/Golang/GoGPUtils/testdata/files/dante.txt"

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	var (
		docBow   DocumentData
		unwanted []string
	)

	if !fileutils.FileExists(filepath) {
		log.Fatal("File " + filepath + " does not exists")
	}

	docBow.DocumentName = filepath
	content, err := ioutil.ReadFile(filepath)

	if err != nil {
		log.Fatal("Unable to read the data for ->" + filepath)
	}

	unwanted = []string{",", ":", ";", ".", "‘", "”", "“", "»", "«", "?", "!"}
	docBow.Bow = StandardizeText(content, true, unwanted)
	log.Println(docBow)
}

// StandardizeText is delegated to generate the BoW for the given data
func StandardizeText(data []byte, toLower bool, toRemove []string) []BoW {
	var (
		// This will contains the text
		text string
		// Text splitted by whitespace
		words []string
		// Total number of words present in the document
		totalWords float64
		// Save the frequencies related to the word
		bow map[string]float64 = make(map[string]float64)
		// Struct for save the BoW and TF
		bowList []BoW
		// Index for insert data into the bowList
		i int
	)

	if toLower {
		log.Println("Lowering text!")

		data = bytes.ToLower(data)
	}
	// Converting []byte in string
	text = string(data)

	// Removing unwanted character/string
	if len(toRemove) > 0 {
		log.Println("Removing the following character from text: [", toRemove, "]")

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

	// Saving the frequencies for each word
	for _, word := range words {
		bow[word]++
		totalWords++
	}

	// Initialize the BoW struct for save the data
	bowList = make([]BoW, len(bow))

	// Save the RAW data into the struct and calculate the TF
	for key, value := range bow {
		bowList[i] = BoW{Word: key, Count: value, TermFrequency: value / totalWords}
		i++
	}
	// Sort the data for the Count
	sort.Slice(bowList, func(i, j int) bool {
		return bowList[i].Count < bowList[j].Count
	})

	return bowList
}
