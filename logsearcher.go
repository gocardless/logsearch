package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mattn/go-isatty"
)

type LogSearcher struct {
	Client     *EsClient
	Query      string
	Period     time.Duration
	NumResults int
	Follow     bool
	idsSeen    map[string]time.Time
	startTime  time.Time
}

func (ls *LogSearcher) Start() {
	ls.startTime = time.Now().Add(-ls.Period)
	ls.idsSeen = make(map[string]time.Time)

	for {
		queryOptions := EsQueryOptions{
			Query:      ls.Query,
			StartTime:  ls.startTime,
			EndTime:    time.Now(),
			NumResults: ls.NumResults,
		}

		resp, err := ls.Client.Search(queryOptions)
		if err != nil {
			log.Fatal(err)
		}

		ls.printResults(resp)

		if !ls.Follow {
			break
		}

		ls.updateStartTime(resp)
		ls.recordIdsSeen(resp)

		time.Sleep(1 * time.Second)
	}
}

func (ls *LogSearcher) printResults(resp *EsResponse) {
	tty := isatty.IsTerminal(os.Stdout.Fd())

	for _, hit := range resp.Hits.Hits {
		if _, ok := ls.idsSeen[hit.Id]; ok {
			continue
		}

		fullMsg, err := json.Marshal(hit.Source)
		if err != nil {
			log.Fatal(err)
		}

		if tty {
			fmt.Printf("\033[34m\033[1m%s\033[0m -- ", hit.Source["@timestamp"])
			fmt.Printf("%s\n\n", string(fullMsg))
		} else {
			fmt.Printf("%s -- %s\n", hit.Source["@timestamp"], string(fullMsg))
		}
	}
}

func (ls *LogSearcher) recordIdsSeen(resp *EsResponse) {
	for id, seenAt := range ls.idsSeen {
		if seenAt.Before(ls.startTime) {
			delete(ls.idsSeen, id)
		}
	}

	for _, hit := range resp.Hits.Hits {
		timestampStr := hit.Source["@timestamp"].(string)
		timestamp, err := time.Parse(time.RFC3339Nano, timestampStr)
		if err == nil {
			ls.idsSeen[hit.Id] = timestamp
		}
	}
}

func (ls *LogSearcher) updateStartTime(resp *EsResponse) {
	newStartTime := time.Now().Add(-(time.Second * 10))
	if newStartTime.After(ls.startTime) {
		ls.startTime = newStartTime
	}
}
