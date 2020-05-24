package index

import (
	"sort"
)

type MatchList struct {
	Matches  int
	Filename string
}

func (inv InvMap) Search(rawQuery string) []MatchList {
	var matchesSlice []MatchList
	var matchesMap = make(map[string]int, 0)

	query := PrepareText(rawQuery)

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
	sort.Slice(matchesSlice, func(i, j int) bool {
		return matchesSlice[i].Matches > matchesSlice[j].Matches
	})

	return matchesSlice
}
