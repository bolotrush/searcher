package index

import (
	"errors"
	"sort"
)

type MatchList struct {
	Matches  int
	Filename string
}

func (inv InvMap) Search(rawQuery string) ([]MatchList, error) {
	var matchesSlice []MatchList
	var matchesMap = make(map[string]int, 0)
	query := PrepareText(rawQuery)
	if len(query) == 0 {
		return nil, errors.New("wrong query")
	}
	for _, word := range query {
		if fileList, ok := inv[word]; ok {
			for _, fileName := range fileList {
				matchesMap[fileName.Filename] += len(fileName.Positions)
			}
		}
	}
	for name, matches := range matchesMap {
		matchesSlice = append(matchesSlice, MatchList{
			Matches:  matches,
			Filename: name,
		})
	}
	if len(matchesSlice) > 0 {
		sort.Slice(matchesSlice, func(i, j int) bool {
			return matchesSlice[i].Matches > matchesSlice[j].Matches
		})
	}
	return matchesSlice, nil
}

func GetDocStrSlice(slice []WordInfo) []string {
	outSlice := make([]string, 0)
	for _, doc := range slice {
		outSlice = append(outSlice, doc.Filename)
	}
	return outSlice
}
