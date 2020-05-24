package index

import (
	"reflect"
	"testing"
)

func getTestInvMap() InvMap {
	newMap := NewInvMap()
	newMap["love"] = []WordInfo{{
		Filename:  "first",
		Positions: []int{0},
	}, {
		Filename:  "second",
		Positions: []int{0},
	}}
	newMap["cats"] = []WordInfo{{
		Filename:  "first",
		Positions: []int{1},
	}}
	return newMap
}

func TestInvMap_InvertIndex(t *testing.T) {
	//in := "love cats."
	filename := "filename"
	in := Token{
		Word:     "love",
		Filename: filename,
		Position: 0,
	}
	expected := NewInvMap()
	expected["love"] = []WordInfo{{
		Filename:  filename,
		Positions: []int{0},
	}}
	actual := NewInvMap()
	actual.AddToken(in)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
	in = Token{
		Word:     "love",
		Filename: filename,
		Position: 2,
	}
	expected["love"] = []WordInfo{{
		Filename:  filename,
		Positions: []int{0, 2},
	}}
	actual.AddToken(in)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
}

func TestInvMap_Search(t *testing.T) {
	in := "love"
	expected := []MatchList{{
		Matches:  1,
		Filename: "first",
	}, {
		Matches:  1,
		Filename: "second",
	}}
	i := getTestInvMap()
	actual := i.Search(in)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
}

func TestIsWordInList(t *testing.T) {
	i := getTestInvMap()
	actual, _ := i.getIndex("love", "second")
	expected := 1
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
	actual, _ = i.getIndex("cats", "first")
	expected = 0
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
}

func TestPrepareText(t *testing.T) {
	in := "I like 254 cats, they are AWESOME!! !"
	expected := []string{"cat", "awesom"}
	actual := PrepareText(in)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
}

func TestPrepareToken(t *testing.T) {
	in := "embarrassment"
	expected := "embarrass"
	actual, _ := prepareToken(in)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
}
