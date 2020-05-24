# Searcher

github.com/bolotrush/searcher implements inverted index to perform full-text search.


## Usage

### Build index

Save inverted index into json file:

    go run main.go file -dir ./examples
    
### Search 

Search by index from file:
    
    go run main.go search -dir ./examples console -q "cat dog" 
    
Search using web-interface. Server address is taken from the environment("config" folder): 

    
    go run main.go search -dir ./examples web
    
To use database connection you need to create PostgreSQL 
with following parameters:
(port=5432 host=localhost user=postgres password=111111 dbname=postgres sslmode=disable)

Tables are located in "database" folder.

Search using web-interface and postgres: 

    go run main.go search -dir ./examples web -db

Search results are sorted by number of found word-tokens.

## Project dependencies

-   [`Stopwords`]("github.com/zoomio/stopwords")
-   [`Snowball`](github.com/kljensen/snowball)
-   [`Cli`](github.com/urfave/cli/v2)
-   [`Zerolog`](github.com/rs/zerolog)
-   [`Env`](github.com/caarlos0/env)
-   [`Go-pg`](github.com/go-pg/pg)