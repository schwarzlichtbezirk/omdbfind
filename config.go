package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// CfgOmdb is OMDb service API calls configuration.
type CfgOmdb struct {
	OmdbHost   string `json:"OMDb-host" yaml:"OMDb-host"`
	ApiKey     string `json:"API-key" yaml:"API-key"`
	FilePath   string `json:"filePath" yaml:"filePath"`
	ThreadsNum int    `json:"threads-num" yaml:"threads-num"`
}

// CfgPrint determines console output parameters.
type CfgPrint struct {
	NoBasicTable bool `json:"no-basic-table" yaml:"no-basic-table"`
	NoOmdbTable  bool `json:"no-OMDb-table" yaml:"no-OMDb-table"`
	TitleLen     int  `json:"title-len" yaml:"title-len"`
	PlotLen      int  `json:"plot-len" yaml:"plot-len"`
}

// CfgSearch is search parameters for local database.
type CfgSearch BasicEntry

// Config is all settings union.
type Config struct {
	MaxRunTime  time.Duration `json:"maxRunTime,omitempty" yaml:"maxRunTime,omitempty"`
	MaxRequests int           `json:"maxRequests,omitempty" yaml:"maxRequests,omitempty"`
	PlotFilter  string        `json:"plotFilter,omitempty" yaml:"plotFilter,omitempty"`
	CfgOmdb     `json:"OMDb-param" yaml:"OMDb-param"`
	CfgSearch   `json:"search-param" yaml:"search-param"`
	CfgPrint    `json:"print-param" yaml:"print-param"`
}

// Instance of settings with program default values.
var cfg = Config{
	MaxRunTime:  0,
	MaxRequests: 100,
	CfgOmdb: CfgOmdb{
		OmdbHost:   "http://omdbapi.com",
		ApiKey:     "124978f0",
		FilePath:   "data-cuted.tsv",
		ThreadsNum: 0,
	},
	CfgPrint: CfgPrint{
		NoBasicTable: true,
		NoOmdbTable:  true,
		TitleLen:     32,
		PlotLen:      32,
	},
}

// compiled binary version, sets by compiler with command
//    go build -ldflags="-X 'main.buildvers=%buildvers%'"
var (
	buildvers string
	_         = buildvers
)

// compiled binary build date, sets by compiler with command
//    go build -ldflags="-X 'main.builddate=%date%'"
var (
	builddate string
	_         = builddate
)

// Register command line flags
func init() {
	flag.StringVar(&cfg.ApiKey, "apiKey", "124978f0", "key for API requests to OMDb service")
	flag.StringVar(&cfg.FilePath, "filePath", "", "absolute path to the inflated title.basics.tsv.gz file")
	flag.DurationVar(&cfg.MaxRunTime, "maxRunTime", 0, "maximum run time of the application. Format is a time.Duration string (for example '1d8h15m30s')")
	flag.IntVar(&cfg.MaxRequests, "maxRequests", 100, "maximum number of requests to send to omdbapi")
	flag.StringVar(&cfg.PlotFilter, "plotFilter", "", "regex pattern to apply to the plot of a film retrieved from OMDb")

	flag.StringVar(&cfg.TitleType, "titleType", "", "filter on titleType column")
	flag.StringVar(&cfg.PrimaryTitle, "primaryTitle", "", "filter on primaryTitle column")
	flag.StringVar(&cfg.OriginalTitle, "originalTitle", "", "filter on originalTitle column")
	flag.IntVar(&cfg.IsAdult, "isAdult", 0, "filter on isAdult column")
	flag.IntVar(&cfg.StartYear, "startYear", 0, "filter on startYear column")
	flag.IntVar(&cfg.EndYear, "endYear", 0, "filter on endYear column")
	flag.IntVar(&cfg.RuntimeMinutes, "runtimeMinutes", 0, "filter on runtimeMinutes column")
	flag.StringVar(&cfg.Genres, "genres", "", "filter on genres column")
}

// ReadYaml reads "data" object from YAML-file with given file name.
func ReadYaml(fname string, data interface{}) (err error) {
	var body []byte
	if body, err = os.ReadFile(filepath.Join(ConfigPath, fname)); err != nil {
		return
	}
	if err = yaml.Unmarshal(body, data); err != nil {
		return
	}
	return
}
