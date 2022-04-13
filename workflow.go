package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"
)

var (
	// context to indicate about service shutdown
	exitctx = context.Background()
	exitfn  = func() {}
	// wait group for all service goroutines
	exitwg sync.WaitGroup
)

var (
	// StartTime, EndTime - start and end time thats program really run.
	StartTime, EndTime time.Time

	// PlotRegexp is compiled regular expression with plot filter.
	PlotRegexp *regexp.Regexp
)

// Init performs global data initialisation.
func Init() {
	log.Printf("version: %s, builton: %s\n", buildvers, builddate)
	log.Println("starts")

	var err error

	// get confiruration path
	if ConfigPath, err = DetectConfigPath(); err != nil {
		log.Fatal(err)
	}
	log.Printf("config path: %s\n", ConfigPath)

	if err = ReadYaml(cfgfile, &cfg); err != nil {
		log.Fatalf("can not read '%s' file: %v\n", cfgfile, err)
	}
	log.Printf("loaded '%s'\n", cfgfile)

	flag.Parse()
}

// WaitExit starts goroutine to wait program termination.
func WaitExit() {
	// checks up command line parameters
	if cfg.MaxRunTime != 0 {
		log.Printf("execution duration is %s\n", cfg.MaxRunTime.String())
	} else {
		log.Printf("no limits on execution time\n")
	}
	if cfg.MaxRequests != 0 {
		log.Printf("maximum %d requests to omdbapi.com\n", cfg.MaxRequests)
	} else {
		log.Printf("no limits on requests numbers\n")
	}
	if cfg.PlotFilter != "" {
		var err error
		if PlotRegexp, err = regexp.Compile(cfg.PlotFilter); err != nil {
			log.Fatal(err)
		}
	}

	// create exit context and wait the break
	if cfg.MaxRunTime != 0 {
		exitctx, exitfn = context.WithTimeout(context.Background(), cfg.MaxRunTime)
	} else {
		exitctx, exitfn = context.WithCancel(context.Background())
	}

	exitwg.Add(1)
	go func() {
		defer exitwg.Done()
		defer exitfn() // make exit signal on function exit

		var sigint = make(chan os.Signal, 1)
		var sigterm = make(chan os.Signal, 1)
		// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGTERM (Ctrl+/)
		// SIGKILL, SIGQUIT will not be caught.
		signal.Notify(sigint, syscall.SIGINT)
		signal.Notify(sigterm, syscall.SIGTERM)
		// Block until we receive our signal.
		select {
		case <-exitctx.Done():
			if errors.Is(exitctx.Err(), context.DeadlineExceeded) {
				log.Println("shutting down by timeout")
			} else if errors.Is(exitctx.Err(), context.Canceled) {
				log.Println("shutting down by cancel")
			} else {
				log.Printf("shutting down by %s", exitctx.Err().Error())
			}
		case <-sigint:
			log.Println("shutting down by break")
		case <-sigterm:
			log.Println("shutting down by process termination")
		}
		signal.Stop(sigint)
		signal.Stop(sigterm)
	}()
}

// Run launches working threads.
func Run() {
	exitwg.Add(1)
	go func() {
		defer exitwg.Done()
		defer exitfn() // let's exit program on func will be complete

		var err error
		// list of database entries thats passes search condition
		var list []BasicEntry

		if list, err = ReadDB(context.Background(), cfg.FilePath); err != nil {
			log.Fatal(err)
		}
		PrintBasic(list)

		// check up on exit
		select {
		case <-exitctx.Done():
			return
		default:
		}

		var res []OmdbEntry
		if res, err = RunPool(list); err != nil {
			log.Fatal(err)
		}
		PrintOmdb(res)
	}()
}

// Done performs graceful network shutdown,
// waits until all server threads will be stopped.
func Done() {
	StartTime = time.Now()
	// wait for exit signal
	<-exitctx.Done()
	// wait until all server threads will be stopped.
	exitwg.Wait()
	EndTime = time.Now()
	var rundur = EndTime.Sub(StartTime)
	log.Printf("time taken: %s\n", rundur.String())
	log.Println("shutting down complete.")
}

func main() {
	Init()
	WaitExit()
	Run()
	Done()
}
