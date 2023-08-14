.DEFAULT_GOAL := build
BINARY_NAME=song-stitch

mod:
	go mod tidy
	go mod vendor

lint: build-ui
	gofmt -s -w cmd/ internal/
	(cd ui && npm run format)
	hadolint Dockerfile

build-ui:
	(cd ui && npm install && npm run build)

run: build-ui
	go run cmd/*.go

watch: build-ui
	gow run cmd/*.go

watch-ui: build-ui
	(cd ui && npm run watch)

run-debug:
	GODEBUG=gctrace=1 go run cmd/*.go

build: lint
	go build -o bin/${BINARY_NAME} cmd/*.go

darwin:
	env GOOS=darwin GOARCH=arm64 go build -o bin/${BINARY_NAME}_darwin_arm64 cmd/*.go
	env GOOS=darwin GOARCH=amd64 go build -o bin/${BINARY_NAME}_darwin_amd64 cmd/*.go
	lipo -create -output bin/${BINARY_NAME}_darwin bin/${BINARY_NAME}_darwin_arm64 bin/${BINARY_NAME}_darwin_amd64

linux-arm64:
	env GOOS=linux GOARCH=arm64 go build -o bin/${BINARY_NAME}_linux_arm64 cmd/*.go

linux-amd64:
	env GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME}_linux_amd64 cmd/*.go

docker-build: lint
	docker-compose build song-stitch

docker-run:
	docker-compose up

clean:
	rm -rf bin/*

gosec:
	gosec -severity medium  ./...

test:
	go clean -testcache
	go test -v ./tests

all: darwin linux-arm64 linux-amd64
