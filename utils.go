package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var durationRegexp = regexp.MustCompile("^(\\d+)\\s*(\\w+)s?$")
var durationUnits = map[string]int64{
	"second": int64(time.Second),
	"minute": int64(time.Minute),
	"hour":   int64(time.Hour),
	"day":    int64(time.Hour) * 24,
	"week":   int64(time.Hour) * 24 * 7,
}

func ParseDuration(duration string) (time.Duration, error) {
	matches := durationRegexp.FindStringSubmatch(duration)
	if len(matches) != 3 {
		return 0, fmt.Errorf("Invalid duration '%s'", duration)
	}

	numUnits, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid duration '%s'", duration)
	}

	unitName := strings.TrimSuffix(matches[2], "s")
	unitLength, ok := durationUnits[unitName]
	if !ok {
		return 0, fmt.Errorf("Invalid duration unit '%s'", unitName)
	}

	return time.Duration(unitLength * numUnits), nil
}
