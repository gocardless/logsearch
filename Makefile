version = 0.2

build:
	go build -ldflags "-X main.version $(version)" .

release:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version $(version)" -o release/logsearch-$(version).linux-amd64/logsearch *.go
	tar -czvf release/logsearch-$(version).linux-amd64.tgz -C release logsearch-$(version).linux-amd64
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version $(version)" -o release/logsearch-$(version).darwin-amd64/logsearch *.go
	tar -czvf release/logsearch-$(version).darwin-amd64.tgz -C release logsearch-$(version).darwin-amd64
