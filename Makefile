all: clean tidy test build run

clean:
	rm -rf zeth
	rm -rf Zeth/
	rm -rf ./**/.ethereum

tidy:
	go mod tidy

test:
	go test -v ./...
.PHONY: test

build:
	go build -o zeth cmd/zeth/main.go

run:
	go run cmd/zeth/*.go

data:
	./scripts/register_nodes.sh
