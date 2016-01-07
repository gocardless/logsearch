package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mattn/go-isatty"
)

const version = "0.1"

var opts struct {
	ElasticsearchUrl string `short:"e" long:"elasticsearch-url" description:"URL for the Elasticsearch instance. Alternatively, use ELASTICSEARCH_URL environment variable."`
	NumResults       int    `short:"n" long:"num-results" description:"Max number of results to return" default:"100"`
	Period           string `short:"p" long:"period" description:"Search time period, e.g. '3 hours', or '1 day'" default:"1 day"`
	Version          bool   `short:"V" long:"version" description:"Show logsearch version"`
}

func printResults(resp *EsResponse) {
	tty := isatty.IsTerminal(os.Stdout.Fd())

	for i, hit := range resp.Hits.Hits {
		fullMsg, err := json.Marshal(hit.Source)
		if err != nil {
			log.Fatal(err)
		}

		if tty {
			if i != 0 {
				fmt.Print("\n")
			}
			fmt.Printf("\033[34m\033[1m%s\033[0m -- ", hit.Source["@timestamp"])
			fmt.Printf("%s\n", string(fullMsg))
		} else {
			fmt.Printf("%s -- %s\n", hit.Source["@timestamp"], string(fullMsg))
		}
	}
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

	queryOptions := EsQueryOptions{
		Query:      args[1],
		StartTime:  time.Now().Add(-searchPeriod),
		EndTime:    time.Now(),
		NumResults: opts.NumResults,
	}

	esUrl := os.Getenv("ELASTICSEARCH_URL")
	if opts.ElasticsearchUrl != "" {
		esUrl = opts.ElasticsearchUrl
	}
	if esUrl == "" {
		fmt.Fprintln(os.Stderr, "Error: missing Elasticsearch URL. Set via ELASTICSEARCH_URL environment variable or the --elasticsearch-url option.")
		os.Exit(1)
	}
	client := EsClient{EsUrl: esUrl}

	resp, err := client.Search(queryOptions)
	if err != nil {
		log.Fatal(err)
	}

	printResults(resp)
}
