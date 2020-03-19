NAME := bqnotify

.PHONY: build
build:
	go build -o ./bin/$(NAME)

.PHONY: build
clean:
	rm bin/$(NAME)

.PHONY: lint
lint:
	golint ./...

.PHONY: fmt
fmt:
	goimports -w .
