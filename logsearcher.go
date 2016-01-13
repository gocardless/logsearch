package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
)

const beginHighlight = "@BEGIN-LOGSEARCH-HIGHLIGHT@"
const endHighlight = "@END-LOGSEARCH-HIGHLIGHT@"
const colorHighlight = "\033[1;7;32m"
const colorReset = "\033[0m"

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

		if tty {
			highlightSourceInline(hit)
		}

		jsonMsgBytes, err := JSONMarshal(hit.Source, true)
		if err != nil {
			log.Fatal(err)
		}
		jsonMsg := string(jsonMsgBytes)

		if tty {
			jsonMsg = strings.Replace(jsonMsg, beginHighlight, colorHighlight, -1)
			jsonMsg = strings.Replace(jsonMsg, endHighlight, colorReset, -1)

			fmt.Printf("\033[34m\033[1m%s\033[0m -- ", hit.Source["@timestamp"])
			fmt.Printf("%s\n\n", jsonMsg)
		} else {
			fmt.Printf("%s -- %s\n", hit.Source["@timestamp"], jsonMsg)
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

func highlightSourceInline(hit EsResponseHit) {
	for fullKey, highlights := range hit.Hightlight {
		if strings.HasSuffix(fullKey, ".raw") {
			continue
		}

		// Deal with fullKeys like headers.accept, etc.
		keyParts := strings.Split(fullKey, ".")
		key := keyParts[0]
		source := hit.Source
		if len(keyParts) > 1 {
			keyParts = keyParts[1:]
			for _, keyPart := range keyParts {
				source = source[key].(map[string]interface{})
				key = keyPart
			}
		}
		if _, ok := source[key]; ok {
			for _, highlight := range highlights {
				source[key] = highlightReplace(source[key], highlight)
			}
		}
	}
}

func highlightReplace(source interface{}, highlight string) interface{} {
	switch source := source.(type) {
	case []interface{}:
		for i, el := range source {
			source[i] = highlightReplace(el, highlight)
		}
		return source
	case map[string]interface{}:
		for key, val := range source {
			source[key] = highlightReplace(val, highlight)
		}
		return source
	case string:
		needle := strings.Replace(highlight, beginHighlight, "", -1)
		needle = strings.Replace(needle, endHighlight, "", -1)
		return strings.Replace(source, needle, highlight, -1)
	}
	return source
}
