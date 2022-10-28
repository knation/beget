
test:
	clear
	go test ./...
.PHONY: test

test.verbose:
	clear
	go test -v ./...
.PHONY: test.verbose

test.cov:
	clear
	go test -coverprofile=coverage.out ./...
.PHONY: test.cov

test.covreport:
	make test.cov
	go tool cover -html=coverage.out
.PHONY: test.coveport

run
	clear
	go build -o build/beget ./main.go

	./build/beget
.PHONY: run