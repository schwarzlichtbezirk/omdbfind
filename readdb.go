package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	ci_tconst = iota
	ci_titleType
	ci_primaryTitle
	ci_originalTitle
	ci_isAdult
	ci_startYear
	ci_endYear
	ci_runtimeMinutes
	ci_genres
)

// BasicEntry is database line format.
type BasicEntry struct {
	TConst         string `json:"tconst,omitempty" yaml:"tconst,omitempty"`
	TitleType      string `json:"titleType,omitempty" yaml:"titleType,omitempty"`
	PrimaryTitle   string `json:"primaryTitle,omitempty" yaml:"primaryTitle,omitempty"`
	OriginalTitle  string `json:"originalTitle,omitempty" yaml:"originalTitle,omitempty"`
	IsAdult        int    `json:"isAdult,omitempty" yaml:"isAdult,omitempty"`
	StartYear      int    `json:"startYear,omitempty" yaml:"startYear,omitempty"`
	EndYear        int    `json:"endYear,omitempty" yaml:"endYear,omitempty"`
	RuntimeMinutes int    `json:"runtimeMinutes,omitempty" yaml:"runtimeMinutes,omitempty"`
	Genres         string `json:"genres,omitempty" yaml:"genres,omitempty"`
}

// Parse converts list of strings to structured data.
func (e *BasicEntry) Parse(record []string) {
	var i64 int64
	e.TConst = record[ci_tconst]
	e.TitleType = record[ci_titleType]
	e.PrimaryTitle = record[ci_primaryTitle]
	e.OriginalTitle = record[ci_originalTitle]
	{
		i64, _ = strconv.ParseInt(record[ci_startYear], 10, 32)
		e.IsAdult = int(i64)
	}
	if record[ci_startYear] != "\\N" {
		i64, _ = strconv.ParseInt(record[ci_startYear], 10, 32)
		e.StartYear = int(i64)
	}
	if record[ci_endYear] != "\\N" {
		i64, _ = strconv.ParseInt(record[ci_endYear], 10, 32)
		e.EndYear = int(i64)
	}
	if record[ci_runtimeMinutes] != "\\N" {
		i64, _ = strconv.ParseInt(record[ci_runtimeMinutes], 10, 32)
		e.RuntimeMinutes = int(i64)
	}
	e.Genres = record[ci_genres]
}

func (e *BasicEntry) Compare(v *BasicEntry) bool {
	if e.TitleType != "" {
		if v.TitleType != e.TitleType {
			return false
		}
	}
	if e.PrimaryTitle != "" {
		if !strings.Contains(strings.ToLower(v.PrimaryTitle), strings.ToLower(e.PrimaryTitle)) {
			return false
		}
	}
	if e.OriginalTitle != "" {
		if !strings.Contains(strings.ToLower(v.OriginalTitle), strings.ToLower(e.OriginalTitle)) {
			return false
		}
	}
	if e.StartYear != 0 {
		if v.StartYear != 0 && v.StartYear != e.StartYear {
			return false
		}
	}
	if e.EndYear != 0 {
		if v.EndYear != 0 && v.EndYear != e.EndYear {
			return false
		}
	}
	if e.RuntimeMinutes != 0 {
		if v.RuntimeMinutes != 0 && v.RuntimeMinutes != e.RuntimeMinutes {
			return false
		}
	}
	if e.Genres != "" {
		var genres = strings.Split(e.Genres, ",")
		var isg = false
		for _, g := range genres {
			if strings.Contains(v.Genres, g) {
				isg = true
			}
		}
		if !isg {
			return false
		}
	}
	return true
}

// ReadDB opens data base with given file name and reads it.
// It applies given from command line filters during reading.
func ReadDB(dbname string) (list []BasicEntry, err error) {
	log.Printf("read file '%s'\n", dbname)

	var f *os.File
	if f, err = os.Open(filepath.Join(ConfigPath, dbname)); err != nil {
		return
	}
	defer f.Close()

	var r = csv.NewReader(f)
	r.Comma = '\t'
	var record []string

	// read and skip header
	if _, err = r.Read(); err != nil {
		return
	}

	var line int
	for {
		if record, err = r.Read(); err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}
		line++

		// format entry and check it
		var v BasicEntry
		v.Parse(record)
		if cfg.Compare(&v) {
			list = append(list, v)
			if cfg.MaxRequests != 0 && len(list) >= cfg.MaxRequests {
				break
			}
		}

		// check up on break by timeout or app termination
		if line%cfg.LineGranulation == 0 {
			select {
			case <-exitctx.Done():
				return
			default:
			}
		}
	}

	return
}

func PrintBasic(list []BasicEntry) {
	if len(list) == 0 {
		log.Printf("no basic entries was found\n")
		return
	}

	log.Printf("founded %d basic entries\n", len(list))
	if !cfg.NoBasicTable {
		fmt.Fprintf(os.Stdout, "IMDB_ID   | Title                            | Year | Genres\n")
		for _, v := range list {
			var t = v.PrimaryTitle
			if len(t) > cfg.TitleLen {
				t = t[:cfg.TitleLen]
			}
			fmt.Fprintf(os.Stdout, "%9s | %-*s | %d | %s\n", v.TConst, cfg.TitleLen, t, v.StartYear, v.Genres)
		}
	}
}
