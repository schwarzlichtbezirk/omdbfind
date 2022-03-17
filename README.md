# OMDb search tool

Utility helps for searching at OMDb service.
Its contains local database with movies descriptions, receives search parameters, gets list of results and queries full description to service for them.

## How to deploy

1. First of all install [Golang](https://go.dev/dl/) of last version. Requires that [GOPATH is set](https://golang.org/doc/code.html#GOPATH). Be sure that `PATH` environment variable contains `%GOPATH%/bin` chunk.

2. Clone this git repo. It will be convenient if you place source code into Golang-folder by original path.

```batch
mkdir %GOPATH%\src\github.com\schwarzlichtbezirk
cd /d %GOPATH%\src\github.com\schwarzlichtbezirk
git clone https://github.com/schwarzlichtbezirk/omdbfind.git
```

3. Download the file: `title.basics.tsv.gz` from [IMDB](https://datasets.imdbws.com/) and unpack it. Place unpacked file `data.tsv` to `github.com/schwarzlichtbezirk/omdbfind/config` path.

4. Compile binaries on explicit OS system (Windows or Linux) by any `build.*` script at `tools` folder. Or build `Dockerfile`.

5. Run testes, or executable, or docker container.

## Configuration

Application can be configured by configuration `omdbfind.yaml` file and by command line parameters. Command line parameters have higher priority then data at `omdbfind.yaml` file, and rewrites its. So, if some parameter is not given at command line, it comes from configuration file. Some parameters sets in code in case if them not provided by any other config.

Command line parameters:

```txt
  -apiKey string
        key for API requests to OMDb service (default "124978f0")
  -endYear int
        filter on endYear column
  -filePath string
        absolute path to the inflated title.basics.tsv.gz file
  -genres string
        filter on genres column
  -isAdult int
        filter on isAdult column
  -maxRequests int
        maximum number of requests to send to omdbapi (default 100)
  -maxRunTime duration
        maximum run time of the application. Format is a time.Duration string (for example '1d8h15m30s')
  -originalTitle string
        filter on originalTitle column
  -plotFilter string
        regex pattern to apply to the plot of a film retrieved from OMDb
  -primaryTitle string
        filter on primaryTitle column
  -runtimeMinutes int
        filter on runtimeMinutes column
  -startYear int
        filter on startYear column
  -titleType string
        filter on titleType column
```

Configuration path with `omdbfind.yaml` file and database is searched at runtime. There is checked `config` path relative to executable, current path, `GOBIN` path, and path relative to debugger.

## Search

Search can be provided by some one parameter, or mix of parameters. The search is greedy, its gets maximum results for given query. If several parameters given, it outputs results with `and` condition.

Search parameters:

* `primaryTitle` - search by primary title. Search is not case sensetive. Can be provided whole title, or a part of title.
* `originalTitle` - search by original title. Same as previous search.
* `isAdult` - search by `isAdult` parameter at database.
* `startYear` - search by start year. If entry have no information about start year in database, its pass this condition on all cases.
* `endYear` - search by end year. Same as previous search.
* `runtimeMinutes` - search by length of movie in minutes. If entry have no information about length in database, its pass this condition on all cases.
* `genres` - search by genre or list of genres divided by comma. Search is case sensetive. If its given a list of genres for search, result passes for each entry if its list of genres contains any genre from condition.

On local database postprocessing there is `plotFilter` filter that points on regular expression for search in `Plot` message.

If it no search parameters provided, all database entries will be given as result.

---
(c) schwarzlichtbezirk, 2022.
