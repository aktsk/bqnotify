NAME := bqnotify

.PHONY: build
build:
	go build -o ./bin/$(NAME)

.PHONY: build
clean:
	rm bin/$(NAME)
