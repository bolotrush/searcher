package files

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"unicode"

	"github.com/bolotrush/searcher/index"

	zl "github.com/rs/zerolog/log"
)

type FileControl struct {
	Mutex    *sync.Mutex
	Wg       *sync.WaitGroup
	dataChan chan index.Token
}

func NewFileControl() FileControl {

	return FileControl{
		Mutex:    &sync.Mutex{},
		Wg:       &sync.WaitGroup{},
		dataChan: make(chan index.Token),
	}

}
func IndexBuild(directory string) (index.InvMap, error) {

	indexMap := index.NewInvMap()
	fc := NewFileControl()

	go fc.listen(&indexMap)

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		go fc.asyncRead(directory, file.Name())
	}

	fc.Wg.Wait()
	close(fc.dataChan)

	return indexMap, nil
}

func (fc *FileControl) listen(indexMap *index.InvMap) {

	for input := range fc.dataChan {
		fc.Wg.Add(1)
		fc.Mutex.Lock()
		indexMap.AddToken(input)
		fc.Mutex.Unlock()
		fc.Wg.Done()
	}
}

func (fc *FileControl) asyncRead(directory string, filename string) {
	text, err := ioutil.ReadFile(directory + "/" + filename)
	if err != nil {
		zl.Err(err).Msg("can not read file")
	}
	words := strings.FieldsFunc(string(text), func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	for position, word := range words {
		token := index.Token{
			Word:     word,
			Filename: strings.TrimRight(filename, ".txt"),
			Position: position,
		}
		fc.dataChan <- token
	}

}

func WriteIndex(indexMap index.InvMap) error {
	file, err := os.Create("out.txt")
	if err != nil {
		return err
	}
	defer closeFile(file)
	indexes, err := json.Marshal(indexMap)
	if err != nil {
		return fmt.Errorf("can not emcode data %w", err)
	}
	if _, err = file.Write(indexes); err != nil {
		return fmt.Errorf("can not write index %w", err)
	}

	return nil
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		zl.Err(err).Msg("can not close file %w")
	}
}
