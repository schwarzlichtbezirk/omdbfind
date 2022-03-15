package main

import (
	"log"
	"testing"
)

func TestQueueDB(t *testing.T) {
	cfg.FilePath = "data-cuted.tsv"
	cfg.MaxRunTime = 0
	cfg.MaxRequests = 0
	cfg.NoBasicTable = true

	log.Println("starts")

	var err error

	// get confiruration path
	if ConfigPath, err = DetectConfigPath(); err != nil {
		t.Fatal(err)
	}
	log.Printf("config path: %s\n", ConfigPath)

	WaitExit()
	Run()
	Done()

	// checkup results
}
