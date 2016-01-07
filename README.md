# logsearch

Search Elasticsearch logs from the command line.

## Installing

If you're using a Mac, install from the Homebrew tap:

```
$ brew tap gocardless/taps
$ brew install logsearch
```

If you're using Linux, download the latest binary from the [releases page][releases].

[releases]: https://github.com/gocardless/logsearch/releases


## Usage

```
$ logsearch -h
Usage:
  logsearch [OPTIONS] QUERY

Application Options:
  -e, --elasticsearch-url= URL for the Elasticsearch instance. Alternatively, use ELASTICSEARCH_URL environment variable.
  -n, --num-results=       Max number of results to return (default: 100)
  -p, --period=            Search time period, e.g. '3 hours', or '1 day' (default: 1 day)

Help Options:
  -h, --help               Show this help message
```

