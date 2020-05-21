package index

import (
	"errors"
	"regexp"
	"strings"

	"github.com/kljensen/snowball"
	"github.com/zoomio/stopwords"
)

type Token struct {
	Word     string
	Filename string
	Position int
}

type Index struct {
	Data     map[string][]WordInfo
	dataChan chan WordInfo
}

type WordInfo struct {
	Filename  string
	Positions []int
}

type InvMap map[string][]WordInfo

func NewInvMap() *InvMap {
	index := make(InvMap)
	return &index
}

func (inv *InvMap) AddToken(token Token) {
	word, err := PrepareToken(token.Word)
	if err != nil {
		//if there's an error just skip word
		return
	}
	if index, ok := inv.inList(word, token.Filename); ok {
		(*inv)[word][index].Positions = append((*inv)[word][index].Positions, token.Position)
	} else {
		structure := WordInfo{
			Filename:  token.Filename,
			Positions: []int{token.Position},
		}
		(*inv)[word] = append((*inv)[word], structure)
	}

}

//func (inv *InvMap) InvertIndex(inputText string, fileName string) {
//	wordList := PrepareText(inputText)
//	for i, word := range wordList {
//		if index, ok := i.inList(word, fileName); ok {
//			(*inv)[word][index].Positions = append((*inv)[word][index].Positions, i)
//		} else {
//			structure := WordInfo{
//				Filename:  fileName,
//				Positions: []int{},
//			}
//			structure.Positions = append(structure.Positions, i)
//			(*i)[word] = append((*i)[word], structure)
//		}
//	}
//}

func (inv InvMap) inList(word string, docId string) (int, bool) {
	for i, ind := range inv[word] {
		if ind.Filename == docId {
			return i, true
		}
	}
	return -1, false
}

var regCompiled = regexp.MustCompile(`[^a-zA-Z_]+`)

func PrepareText(in string) []string {
	tokens := clean(regCompiled.Split(in, -1))
	return tokens
}

func clean(inputWords []string) []string {
	cleanWords := make([]string, 0)
	for _, word := range inputWords {
		if stopwords.IsStopWord(word) {
			continue
		}
		word = strings.ToLower(word)
		cleanWords = append(cleanWords, word)
	}
	return cleanWords
}

func PrepareToken(word string) (string, error) {
	if stopwords.IsStopWord(word) {
		return "", errors.New("stop word")
	}
	return snowball.Stem(word, "english", true)
}
