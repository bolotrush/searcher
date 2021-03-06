package db

import (
	"fmt"

	zl "github.com/rs/zerolog/log"

	"github.com/bolotrush/searcher/index"
	"github.com/go-pg/pg"
)

type Word struct {
	ID   int    `sql:"id,pk"`
	Word string `sql:"word"`
}
type File struct {
	ID   int    `sql:"id,pk"`
	Name string `sql:"name"`
}
type Occurrence struct {
	ID       int `sql:"id,pk"`
	WordID   int `sql:"word_id"`
	FileID   int `sql:"file_id"`
	Position int `sql:"position"`
}
type Base struct {
	pg *pg.DB
}

func NewDb(config string) (*Base, error) {
	options, err := pg.ParseURL(config)
	if err != nil {
		return &Base{}, err
	}
	zl.Debug().
		Str("database", options.Database).
		Str("user", options.User).
		Str("address", options.Addr).
		Msg("db is created")
	return &Base{
		pg: pg.Connect(options),
	}, nil
}

func (b *Base) WriteIndex(index index.InvMap) error {
	if err := b.clearTables(); err != nil {
		return fmt.Errorf("can't clear data in db %w", err)
	}
	if err := b.addIndex(index); err != nil {
		return fmt.Errorf("can't add indexes into db %w", err)
	}
	return nil
}

func (b *Base) Search(rawQuery string) []index.MatchList {
	query := index.PrepareText(rawQuery)
	var result []index.MatchList
	var occ Occurrence
	err := b.pg.Model(&occ).
		ColumnExpr("files.name AS filename").
		ColumnExpr("count(position) AS matches").
		Join("JOIN words ON words.id = word_id").
		Join("JOIN files ON files.id = file_id").
		WhereIn("words.word IN (?)", query).
		Group("files.name").
		Order("matches DESC").
		Select(&result)
	if err != nil {
		zl.Err(err).Msgf("cant get results: %w", err)
	}
	return result
}

func (b *Base) Close() {
	err := b.pg.Close()
	if err != nil {
		zl.Err(err).Msg("can't close db connection")
	}
}

func (b *Base) clearTables() error {
	if _, err := b.pg.Exec("TRUNCATE words, files, occurrences;"); err != nil {
		return err
	}
	return nil
}

func (b *Base) addIndex(index index.InvMap) error {
	for token, properties := range index {
		token := Word{
			Word: token,
		}
		err := b.pg.Insert(&token)
		if err != nil {
			return err
		}
		for _, property := range properties {
			file := File{
				Name: property.Filename,
			}
			err := b.pg.Insert(&file)
			if err != nil {
				return err
			}
			var occurrences []Occurrence
			for _, position := range property.Positions {
				occurrence := Occurrence{
					WordID:   token.ID,
					FileID:   file.ID,
					Position: position,
				}
				occurrences = append(occurrences, occurrence)
			}
			err = b.pg.Insert(&occurrences)
			if err != nil {
				return err
			}

		}
	}
	return nil
}

func (b *Base) getTokenID(query []string) ([]int, error) {
	var words Word
	var wordIds []int
	err := b.pg.Model(&words).
		Column("id").
		WhereIn("word IN (?)", query).
		Select(&wordIds)
	if err != nil {
		return nil, fmt.Errorf("can't get token from db %w", err)
	}
	return wordIds, nil
}
