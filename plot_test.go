package main

import (
	"context"
	"log"
	"testing"
)

func doOmdb(t *testing.T, idtt []string) {
	cfg.FilePath = "data-cuted.tsv"
	cfg.MaxRunTime = 0
	cfg.MaxRequests = 0

	log.Println("starts")

	var err error

	// get confiruration path
	if ConfigPath, err = DetectConfigPath(); err != nil {
		t.Fatal(err)
	}
	log.Printf("config path: %s\n", ConfigPath)

	WaitExit()

	var res []OmdbEntry
	exitwg.Add(1)
	go func() {
		defer exitwg.Done()
		defer exitfn() // let's exit program on func will be complete
		// list of database entries thats passes search condition
		var list []BasicEntry

		if list, err = ReadDB(context.Background(), cfg.FilePath); err != nil {
			log.Fatal(err)
		}
		PrintBasic(list)

		if res, err = RunPool(list); err != nil {
			log.Fatal(err)
		}
		PrintOmdb(res)
	}()

	Done()

	// checkup results
	if len(res) != len(idtt) {
		t.Fatalf("expected %d entries, got %d", len(idtt), len(res))
	}
	for i, id := range idtt {
		if res[i].ImdbID != id {
			t.Fatalf("wrong search result on entry #%d, wanted '%s', got '%s'", i, id, res[i].ImdbID)
		}
	}
}

// TestPlot1 tests with regexp plot filter "^Two".
func TestPlot1(t *testing.T) {
	cfg.PlotFilter = "^Two"

	doOmdb(t, []string{
		"tt0000017", "tt0000022", "tt0000026",
	})
}

// TestPlot2 tests with regexp plot filter "Two".
func TestPlot2(t *testing.T) {
	cfg.PlotFilter = "Two"

	doOmdb(t, []string{
		"tt0000016", "tt0000017", "tt0000022", "tt0000026",
	})
}
