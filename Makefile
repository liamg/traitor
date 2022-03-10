

.PHONY: build
build:
	CGO_ENABLED=0 go build ./cmd/traitor1

.PHONY: pack
pack:
	go run ./cmd/pack

.PHONY: install
install:
	

.PHONY: test
test:
	go test ./... -race -cover
