package main

import (
	"errors"
	"fmt"
	"github.com/bolotrush/searcher/models/db"
	"os"

	"github.com/bolotrush/searcher/config"
	"github.com/bolotrush/searcher/models/files"
	"github.com/bolotrush/searcher/web"
	"github.com/rs/zerolog"

	"github.com/urfave/cli/v2"

	zl "github.com/rs/zerolog/log"
)

var cfg config.Config

func main() {
	var err error
	cfg, err = config.Load()
	if err != nil {
		zl.Fatal().Err(err).Msg("can not load configs")
	}
	loglevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		zl.Fatal().Err(err).Msg("can not parse log level")
	}
	zerolog.SetGlobalLevel(loglevel)

	app := &cli.App{
		Name:  "Searcher",
		Usage: "The app searches docs using inverted index and find the best match",
	}
	dirFlag := &cli.StringFlag{
		Name:     "dir",
		Usage:    "Path to files directory",
		Required: true,
	}
	app.Commands = []*cli.Command{
		{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "Saves index to file",
			Flags:   []cli.Flag{dirFlag},
			Action:  indexFile,
		}, {
			Name:  "search",
			Usage: "Reads query and search files",
			Flags: []cli.Flag{dirFlag},
			Subcommands: []*cli.Command{
				{
					Name:  "console",
					Usage: "Searches index in console",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "query",
							Aliases:  []string{"q"},
							Usage:    "Searches query in console",
							Required: true,
						},
					},
					Action: searchConsole,
				}, {
					Name:  "web",
					Usage: "Creates web server for search using http",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:  "db",
							Usage: "Uses PostgreSQL for data",
						},
					},
					Action: searchWeb,
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		zl.Fatal().Err(err).Msg("app run error")
	}
}

func indexFile(c *cli.Context) error {
	inPath := c.String("dir")
	indexMap, err := files.IndexBuild(inPath)
	if err != nil {
		return fmt.Errorf("can't create index %w", err)
	}
	return files.WriteIndex(indexMap)
}

func searchConsole(c *cli.Context) error {
	path := c.String("dir")
	query := c.String("query")

	indexMap, err := files.IndexBuild(path)
	if err != nil {
		return fmt.Errorf("can't create index %w", err)
	}

	matches := indexMap.Search(query)
	if len(matches) > 0 {
		for i, match := range matches {
			fmt.Printf("%d) %s: matches - %d\n", i+1, match.Filename, match.Matches)
		}
	} else {
		fmt.Println("There's no results :(")
	}
	return nil
}

func searchWeb(c *cli.Context) error {
	path := c.String("dir")
	if len(path) == 0 {
		return errors.New("path to files is empty")
	}
	indexMap, err := files.IndexBuild(path)
	if err != nil {
		return fmt.Errorf("can't create index %w", err)
	}

	var searcher web.SearchFunc
	if c.Bool("db") {
		base, err := db.NewDb(cfg.PgSQL)
		if err != nil {
			return fmt.Errorf("error on creating db %w", err)
		}
		defer base.Close()

		if err := base.WriteIndex(indexMap); err != nil {
			return fmt.Errorf("error on db index writing %w", err)
		}
		searcher = base.Search
	} else {
		searcher = indexMap.Search
	}

	server, err := web.NewServer(cfg.Listen, searcher)
	if err != nil {
		return fmt.Errorf("can't create server: %w", err)
	}
	return server.Run()
}
