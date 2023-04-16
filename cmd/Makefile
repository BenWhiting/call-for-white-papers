.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"  -buildvcs=false  -o build/papers-please ${MAIN_FILE}

.PHONY: run
run:
	./build/papers-please