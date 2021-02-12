

build:
	CGO_ENABLED=0 go build ./cmd/traitor

install:
	CGO_ENABLED=0 go install -ldflags "-X github.com/liamg/traitor/version.Version=`git describe --tags`" ./cmd/traitor

test:
	go test ./... -race -cover
