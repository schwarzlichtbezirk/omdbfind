package main

import (
	"log"
	"testing"
)

func doBasic(t *testing.T, idtt []string) {
	cfg.FilePath = "data-cuted.tsv"

	log.Println("starts")

	var err error
	var list []BasicEntry

	// get confiruration path
	if ConfigPath, err = DetectConfigPath(); err != nil {
		t.Fatal(err)
	}
	log.Printf("config path: %s\n", ConfigPath)

	WaitExit()
	exitwg.Add(1)
	go func() {
		defer exitwg.Done()
		defer exitfn() // let's exit program on func will be complete

		var err error
		if list, err = ReadDB(cfg.FilePath); err != nil {
			log.Fatal(err)
		}
		PrintBasic(list)
	}()
	Done()

	// checkup results
	if len(list) != len(idtt) {
		t.Fatalf("expected %d entries, got %d", len(idtt), len(list))
	}
	for i, id := range idtt {
		if list[i].TConst != id {
			t.Fatalf("wrong search result on entry #%d, wanted '%s', got '%s'", i, id, list[i].TConst)
		}
	}
}

// Test search by the whole title.
// Returns any entry that have equal title.
// Search is not case sensitive.
func TestTitle(t *testing.T) {
	// setup search filter
	cfg.PrimaryTitle = "watering the flowers"

	doBasic(t, []string{
		"tt0000035",
	})
}

// Test search by the part of title.
// Returns any entry thats title contains given string.
// Search is not case sensitive.
func TestWord(t *testing.T) {
	// setup search filter
	cfg.PrimaryTitle = "clown"

	doBasic(t, []string{
		"tt0000002", "tt0000019",
	})
}

// Test search by the start year.
// Returns any entries with equal start year.
func TestYear(t *testing.T) {
	// setup search filter
	cfg.StartYear = 1894

	doBasic(t, []string{
		"tt0000001", "tt0000006", "tt0000007", "tt0000008", "tt0000009", "tt0000015",
	})
}

// Test search by the single genre.
// Returns any entries that contains given genre in the its genres list.
// Genres are case sensitive.
func TestGenre(t *testing.T) {
	// setup search filter
	cfg.Genres = "Comedy"

	doBasic(t, []string{
		"tt0000003", "tt0000005", "tt0000014", "tt0000019", "tt0000033", "tt0000035",
	})
}

// Test search by the list of genres.
// Returns any entries that have some genre from the list.
// Genres are case sensitive.
func TestGenres(t *testing.T) {
	// setup search filter
	cfg.Genres = "Comedy,Drama"

	doBasic(t, []string{
		"tt0000003", "tt0000005", "tt0000014", "tt0000019", "tt0000033", "tt0000035", "tt0000036",
	})
}

// Test search by the mix of conditions.
func TestMix(t *testing.T) {
	// setup search filter
	cfg.StartYear = 1894
	cfg.Genres = "Animation,Documentary"

	doBasic(t, []string{
		"tt0000001", "tt0000008", "tt0000015",
	})
}
