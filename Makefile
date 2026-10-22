BIN := bin

.PHONY: build test clean
build:
	go build -o $(BIN)/mergevcf ./cmd/mergevcf
	go build -o $(BIN)/cleanvcf ./cmd/cleanvcf
	go build -o $(BIN)/showvcf ./cmd/showvcf

test:
	go test ./...

clean:
	rm -rf $(BIN)
