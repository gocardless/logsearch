package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

var version string

var opts struct {
	ElasticsearchUrl string `short:"e" long:"elasticsearch-url" description:"URL for the Elasticsearch instance. Alternatively, use ELASTICSEARCH_URL environment variable."`
	Follow           bool   `short:"f" long:"follow" description:"Show new data as it becomes available, like tail -f"`
	NumResults       int    `short:"n" long:"num-results" description:"Max number of results to return" default:"100"`
	Period           string `short:"p" long:"period" description:"Search time period, e.g. '3 hours', or '1 day'" default:"1 day"`
	Version          bool   `short:"V" long:"version" description:"Show logsearch version"`
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.Usage = "[OPTIONS] QUERY"

	args, err := parser.ParseArgs(os.Args)
	if err != nil {
		os.Exit(1)
	}

	if opts.Version {
		fmt.Printf("logsearch v%s\n", version)
		return
	}

	searchPeriod, err := ParseDuration(opts.Period)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(args) != 2 {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	esUrl := os.Getenv("ELASTICSEARCH_URL")
	if opts.ElasticsearchUrl != "" {
		esUrl = opts.ElasticsearchUrl
	}
	if esUrl == "" {
		fmt.Fprintln(os.Stderr, "Error: missing Elasticsearch URL. Set via ELASTICSEARCH_URL environment variable or the --elasticsearch-url option.")
		os.Exit(1)
	}
	client := &EsClient{EsUrl: esUrl}

	searcher := LogSearcher{
		Client:     client,
		Query:      args[1],
		Period:     searchPeriod,
		NumResults: opts.NumResults,
		Follow:     opts.Follow,
	}
	searcher.Start()
}
