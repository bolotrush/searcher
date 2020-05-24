/*
Package index implements inverted index with functions to index new documents, to search over the built
index. Inverted index uses in-memory or database engine.
*/
package index

import (
	"errors"
	"github.com/kljensen/snowball"
	"github.com/zoomio/stopwords"
	"strings"
	"unicode"
)

//Token is the struct contains information about every word in file.
//Instances of this struct are needed for creating inverted index
type Token struct {
	Word     string
	Filename string
	Position int
}

//InvMap is the signature type that contains inverted index. Keys of this map are words in text file
type InvMap map[string][]WordInfo

//WordInfo is the value of InvMap that contains filename and slice of word's positions in text file
type WordInfo struct {
	Filename  string
	Positions []int
}

//NewInvMap creates new object of InvMap type
func NewInvMap() InvMap {
	index := make(InvMap)
	return index
}

//AddToken adds tokens into an InvMap
func (inv *InvMap) AddToken(token Token) {

	if index, ok := inv.getIndex(token.Word, token.Filename); ok {
		(*inv)[token.Word][index].Positions = append((*inv)[token.Word][index].Positions, token.Position)
	} else {
		structure := WordInfo{
			Filename:  token.Filename,
			Positions: []int{token.Position},
		}
		(*inv)[token.Word] = append((*inv)[token.Word], structure)
	}
}

const MinWordLen = 2

//PrepareToken delete stop words and short words
func PrepareText(text string) []string {
	var prepared []string
	tokens := strings.FieldsFunc(string(text), func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	for _, token := range tokens {
		cleanToken, err := prepareToken(token)
		if err != nil {
			//if there's an error while stemming word just skip it
			continue
		}
		prepared = append(prepared, cleanToken)
	}
	return prepared
}

func prepareToken(word string) (string, error) {
	if stopwords.IsStopWord(word) || len(word) < MinWordLen {
		return "", errors.New("no need this word")
	}
	return snowball.Stem(word, "english", true)
}

func (inv InvMap) getIndex(word string, docId string) (int, bool) {
	for i, ind := range inv[word] {
		if ind.Filename == docId {
			return i, true
		}
	}
	return -1, false
}
