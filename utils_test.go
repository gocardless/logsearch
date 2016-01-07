package main

import (
	"testing"
	"time"
)

func TestParseSeconds(t *testing.T) {
	assertParseFailure(t, "seconds")
	assertParseResult(t, "1 second", time.Second)
	assertParseResult(t, "2 seconds", time.Second*2)
}

func TestParseMinutes(t *testing.T) {
	assertParseFailure(t, "minutes")
	assertParseResult(t, "1 minute", time.Minute)
	assertParseResult(t, "2 minutes", time.Minute*2)
}

func TestParseHours(t *testing.T) {
	assertParseFailure(t, "hours")
	assertParseResult(t, "1 hour", time.Hour)
	assertParseResult(t, "2 hours", time.Hour*2)
}

func TestParseDays(t *testing.T) {
	assertParseFailure(t, "days")
	assertParseResult(t, "1 day", time.Hour*24)
	assertParseResult(t, "2 days", time.Hour*48)
}

func TestParseWeeks(t *testing.T) {
	assertParseFailure(t, "weeks")
	assertParseResult(t, "1 week", time.Hour*24*7)
	assertParseResult(t, "2 weeks", time.Hour*24*14)
}

func assertParseFailure(t *testing.T, str string) {
	_, err := ParseDuration(str)
	if err == nil {
		t.Fatalf("expected error parsing 'seconds', got none")
	}
}

func assertParseResult(t *testing.T, str string, duration time.Duration) {
	parsedDuration, err := ParseDuration(str)
	if err != nil {
		t.Fatalf("unexpected error parsing '%s': %s", str, err)
	}

	if duration != parsedDuration {
		t.Fatalf("expected %d, got %d", duration, parsedDuration)
	}
}
