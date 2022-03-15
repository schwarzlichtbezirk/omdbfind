package main

import (
	"log"
	"testing"
)

// doQueueDB runs benchmark with given number of threads
// of queries to OMDb service.
func doQueueDB(b *testing.B, thrnum int) {
	cfg.FilePath = "data-cuted.tsv"
	cfg.MaxRunTime = 0
	cfg.MaxRequests = 0
	cfg.NoBasicTable = true
	cfg.NoOmdbTable = true
	cfg.ThreadsNum = thrnum

	log.Println("starts")

	var err error

	// get confiruration path
	if ConfigPath, err = DetectConfigPath(); err != nil {
		log.Fatal(err)
	}
	log.Printf("config path: %s\n", ConfigPath)

	// list of database entries thats passes search condition
	var list []BasicEntry

	if list, err = ReadDB(cfg.FilePath); err != nil {
		log.Fatal(err)
	}
	PrintBasic(list)

	for i := 0; i < b.N; i++ {
		log.Printf("iter #%d of %d", i, b.N)
		WaitExit()

		exitwg.Add(1)
		go func() {
			defer exitwg.Done()
			defer exitfn() // let's exit program on func will be complete

			var res = RunPool(list)
			PrintOmdb(res)
		}()

		Done()
	}
}

// BenchmarkThrNum1 runs benchmark with 1 thread of queries to OMDb service.
func BenchmarkThrNum1(b *testing.B) {
	doQueueDB(b, 1)
}

// BenchmarkThrNum2 runs benchmark with 2 threads of queries to OMDb service.
func BenchmarkThrNum2(b *testing.B) {
	doQueueDB(b, 2)
}

// BenchmarkThrNum4 runs benchmark with 4 threads of queries to OMDb service.
func BenchmarkThrNum4(b *testing.B) {
	doQueueDB(b, 4)
}

// BenchmarkThrNum0 runs benchmark with threads on each CPUs of queries to OMDb service.
func BenchmarkThrNum0(b *testing.B) {
	doQueueDB(b, 0)
}
