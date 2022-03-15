package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Rating struct {
	Source string
	Value  string
}

type OmdbEntry struct {
	Title      string   `json:",omitempty" yaml:",omitempty"`
	Year       string   `json:",omitempty" yaml:",omitempty"`
	Rated      string   `json:",omitempty" yaml:",omitempty"`
	Released   string   `json:",omitempty" yaml:",omitempty"`
	Runtime    string   `json:",omitempty" yaml:",omitempty"`
	Genre      string   `json:",omitempty" yaml:",omitempty"`
	Director   string   `json:",omitempty" yaml:",omitempty"`
	Writer     string   `json:",omitempty" yaml:",omitempty"`
	Actors     string   `json:",omitempty" yaml:",omitempty"`
	Plot       string   `json:",omitempty" yaml:",omitempty"`
	Language   string   `json:",omitempty" yaml:",omitempty"`
	Country    string   `json:",omitempty" yaml:",omitempty"`
	Awards     string   `json:",omitempty" yaml:",omitempty"`
	Poster     string   `json:",omitempty" yaml:",omitempty"`
	Ratings    []Rating `json:",omitempty" yaml:",omitempty"`
	Metascore  string   `json:",omitempty" yaml:",omitempty"`
	ImdbRating string   `json:"imdbRating,omitempty" yaml:"imdbRating,omitempty"`
	ImdbVotes  string   `json:"imdbVotes,omitempty" yaml:"imdbVotes,omitempty"`
	ImdbID     string   `json:"imdbID,omitempty" yaml:"imdbID,omitempty"`
	Type       string   `json:",omitempty" yaml:",omitempty"`
	DVD        string   `json:",omitempty" yaml:",omitempty"`
	BoxOffice  string   `json:",omitempty" yaml:",omitempty"`
	Production string   `json:",omitempty" yaml:",omitempty"`
	Website    string   `json:",omitempty" yaml:",omitempty"`
	Response   string   `json:",omitempty" yaml:",omitempty"`
}

// QueryDB is gets data for each entry step by step.
// Its a part of multithreaded pipeline.
func QueryDB(list []BasicEntry) (res []OmdbEntry, err error) {
	for _, be := range list {
		var requri = fmt.Sprintf("%s/?i=%s&apikey=%s", cfg.OmdbHost, be.TConst, cfg.ApiKey)
		var resp *http.Response
		if resp, err = http.Get(requri); err != nil {
			return
		}
		var body []byte
		body, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return
		}

		var oe OmdbEntry
		err = json.Unmarshal(body, &oe)
		res = append(res, oe)

		// check up on exit
		select {
		case <-exitctx.Done():
			return
		default:
		}
	}
	return
}

func PrintOmdb(list []OmdbEntry) {
	if len(list) == 0 {
		log.Printf("no OMDb entries was found\n")
		return
	}

	log.Printf("founded %d OMDb entries\n", len(list))
	if !cfg.NoBasicTable {
		fmt.Fprintf(os.Stdout, "IMDB_ID   | Title                            | Plot\n")
		for _, v := range list {
			var t = v.Title
			if len(t) > cfg.TitleLen {
				t = t[:cfg.TitleLen]
			}
			var p = v.Plot
			if len(p) > cfg.PlotLen {
				p = p[:cfg.PlotLen]
			}
			fmt.Fprintf(os.Stdout, "%9s | %-*s | %-*s\n", v.ImdbID, cfg.TitleLen, t, cfg.PlotLen, p)
		}
	}
}
