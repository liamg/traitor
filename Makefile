

.PHONY: build
build:
	CGO_ENABLED=0 go build ./cmd/traitor1



.PHONY: install
install:
	

.PHONY: test
test:
	go test ./... -race -cover
