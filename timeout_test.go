package main

import (
	"log"
	"testing"
	"time"
)

// timeout_limit is predefined time limit of program execution.
const timeout_limit = 1000 * time.Microsecond

// timeout_diff is maximum difference between
// real program running timeout and given limit.
// it should be relative to OS timer granulation.
const timeout_diff = 5 * time.Millisecond

func TestTimeout(t *testing.T) {
	cfg.FilePath = "data.tsv"
	cfg.MaxRunTime = timeout_limit
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
	var rundur = EndTime.Sub(StartTime)
	if rundur < cfg.MaxRunTime {
		t.Fatal("really running time is less then given")
	}
	if rundur-cfg.MaxRunTime > timeout_diff {
		t.Fatal("running time exceeds given limit")
	}
}
