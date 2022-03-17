package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	// ErrColumnsNum returns if given row to table print have
	// columns number not equal to number columns at header.
	ErrColumnsNum = errors.New("row contains wrong unexpected number of columns")
)

// CellInfo describes structure field output as table cell.
type CellInfo struct {
	Value  string `json:"value" yaml:"value"`
	Limit  int    `json:"limit,omitempty" yaml:"limit,omitempty"`
	Margin int    `json:"margin,omitempty" yaml:"margin,omitempty"`
}

// Print prints single cell of table.
func (ci *CellInfo) Print(w io.Writer) (n int, err error) {
	var val, lim = ci.Value, ci.Limit
	if lim > 0 && len(val) > lim {
		val = val[:lim]
	}
	var sep = strings.Repeat(" ", ci.Margin)
	return fmt.Fprintf(w, "%s%-*s%s", sep, lim, val, sep)
}

// TableInfo describes whole table output.
type TableInfo struct {
	Header []CellInfo `json:"header" yaml:"header"`
	Border bool       `json:"border,omitempty" yaml:"border,omitempty"`
	Grid   bool       `json:"grid,omitempty" yaml:"grid,omitempty"`
}

// PrintRow prints whole table row.
func (ti *TableInfo) PrintRow(w io.Writer, row []string) (n int, err error) {
	if len(row) != len(ti.Header) {
		err = ErrColumnsNum
		return
	}
	var ni int
	for i, val := range row {
		if (i == 0 && ti.Border) || (i != 0 && ti.Grid) {
			ni, err = fmt.Fprintf(w, "|")
			n += ni
			if err != nil {
				return
			}
		}
		var ci = ti.Header[i]
		ci.Value = val
		ni, err = ci.Print(w)
		n += ni
		if err != nil {
			return
		}
	}
	if ti.Border {
		ni, err = fmt.Fprintf(w, "|")
		n += ni
		if err != nil {
			return
		}
	}
	ni, err = fmt.Fprintf(w, "\n")
	n += ni
	return
}

// PrintTable prints content of whole table.
func PrintTable[T Fielder](w io.Writer, ti *TableInfo, list []T) (n int, err error) {
	var ni int
	var row = make([]string, len(ti.Header))
	for i, ci := range ti.Header {
		row[i] = ci.Value
	}
	ni, err = ti.PrintRow(w, row)
	n += ni
	if err != nil {
		return
	}
	for _, t := range list {
		for i, ci := range ti.Header {
			row[i], _ = t.Field(ci.Value)
		}
		ni, err = ti.PrintRow(w, row)
		n += ni
		if err != nil {
			return
		}
	}
	return
}

// Fielder helps to call "Field" method to get
// namded field value of structure.
type Fielder interface {
	Field(name string) (fld string, ok bool)
}

// Field returns structure field content by it's name.
func (be BasicEntry) Field(name string) (fld string, ok bool) {
	ok = true
	switch name {
	case "IMDB_ID", "tconst":
		fld = be.TConst
	case "titleType":
		fld = be.TitleType
	case "Title", "primaryTitle":
		fld = be.PrimaryTitle
	case "originalTitle":
		fld = be.OriginalTitle
	case "isAdult":
		fld = strconv.Itoa(be.IsAdult)
	case "Year", "startYear":
		fld = strconv.Itoa(be.StartYear)
	case "endYear":
		fld = strconv.Itoa(be.EndYear)
	case "Length", "runtimeMinutes":
		fld = strconv.Itoa(be.RuntimeMinutes)
	case "Genres", "genres":
		fld = be.Genres
	default:
		ok = false
	}
	return
}

// Field returns structure field content by it's name.
func (oe OmdbEntry) Field(name string) (fld string, ok bool) {
	ok = true
	switch name {
	case "Title":
		fld = oe.Title
	case "Year":
		fld = oe.Year
	case "Rated":
		fld = oe.Rated
	case "Released":
		fld = oe.Released
	case "Runtime":
		fld = oe.Runtime
	case "Genre":
		fld = oe.Genre
	case "Director":
		fld = oe.Director
	case "Writer":
		fld = oe.Writer
	case "Actors":
		fld = oe.Actors
	case "Plot":
		fld = oe.Plot
	case "Language":
		fld = oe.Language
	case "Country":
		fld = oe.Country
	case "Awards":
		fld = oe.Awards
	case "Poster":
		fld = oe.Poster
	case "Metascore":
		fld = oe.Metascore
	case "ImdbRating", "imdbRating":
		fld = oe.ImdbRating
	case "ImdbVotes", "imdbVotes":
		fld = oe.ImdbVotes
	case "ImdbID", "imdbID", "IMDB_ID", "tconst":
		fld = oe.ImdbID
	case "Type":
		fld = oe.Type
	case "DVD":
		fld = oe.DVD
	case "BoxOffice":
		fld = oe.BoxOffice
	case "Production":
		fld = oe.Production
	case "Website":
		fld = oe.Website
	case "Response":
		fld = oe.Response
	default:
		ok = false
	}
	return
}

// PrintBasic prints formated table of BasicEntry content to console.
func PrintBasic(list []BasicEntry) {
	if len(list) == 0 {
		log.Printf("no basic entries was found\n")
		return
	}

	log.Printf("founded %d basic entries\n", len(list))
	if !cfg.NoBasicTable {
		PrintTable(os.Stdout, &cfg.BasicTableInfo, list)
	}
}

// PrintOmdb prints formated table of OmdbEntry content to console.
func PrintOmdb(list []OmdbEntry) {
	if len(list) == 0 {
		log.Printf("no OMDb entries was found\n")
		return
	}

	log.Printf("founded %d OMDb entries\n", len(list))
	if !cfg.NoOmdbTable {
		PrintTable(os.Stdout, &cfg.OmdbTableInfo, list)
	}
}
