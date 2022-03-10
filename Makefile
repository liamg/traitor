

.PHONY: build
build:
	CGO_ENABLED=0 go build ./cmd/traitor

.PHONY: pack
pack:
	go run ./cmd/pack

.PHONY: install
install:
	CGO_ENABLED=0 go install -ldflags "-X github.com/liamg/traitor/version.Version=`git describe --tags`" ./cmd/traitor

.PHONY: test
test:
	go test ./... -race -cover
