package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"sync"
)

// Rating describes rating data returned by queries to OMDb service.
type Rating struct {
	Source string
	Value  string
}

// OmdbEntry is movies description returned by queries to OMDb service.
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
func QueryDB(ctx context.Context, list []BasicEntry) (res []OmdbEntry, err error) {
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
		if err = json.Unmarshal(body, &oe); err != nil {
			return
		}
		oe.ImdbID = be.TConst

		var add = true
		if PlotRegexp != nil {
			add = PlotRegexp.MatchString(oe.Plot)
		}
		if add {
			res = append(res, oe)
		}

		// check up on exit
		select {
		case <-ctx.Done():
			return
		case <-exitctx.Done():
			return
		default:
		}
	}
	return
}

// RunPool starts a pool of threads to get data
// from service for excerpted entries.
func RunPool(list []BasicEntry) (res []OmdbEntry, err error) {
	var ll = len(list)
	var tn = cfg.ThreadsNum
	if tn <= 0 {
		tn = runtime.NumCPU()
	}
	if tn > ll {
		tn = ll
	}
	var tres = make([][]OmdbEntry, tn)
	log.Printf("starts %d threads to OMDb service\n", tn)

	var cherr = make(chan error, tn)
	var ctx, ctxfn = context.WithCancel(context.Background())
	go func() {
		defer ctxfn()
		select {
		case err = <-cherr:
			return
		case <-ctx.Done():
			return
		}
	}()

	var wg sync.WaitGroup
	wg.Add(tn)
	for i := 0; i < tn; i++ {
		var i = i // localize
		go func() {
			defer wg.Done()
			var in = ll / tn
			var sl = in
			if i == tn-1 {
				sl += ll % tn
			}
			if tres[i], err = QueryDB(ctx, list[i*in:i*in+sl]); err != nil {
				select {
				case <-ctx.Done():
					return
				case cherr <- err:
					return
				}
			}
		}()
	}
	wg.Wait()
	close(cherr) // err will be guaranteed with expected value
	<-ctx.Done() // wait until error goroutine exit
	if err != nil {
		return
	}

	for _, ires := range tres {
		res = append(res, ires...)
	}
	return
}
