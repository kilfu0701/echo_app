NAME     = main
PID      = /tmp/echo_app.pid
GO_FILES = src/*.go
APP      = ./bin/main

bin/$(NAME):
	$(eval sha := $(shell git rev-parse HEAD))
	go build -i -ldflags="-s -w -X main.gitHash=$(sha)" -o bin/$(NAME) $(GO_FILES)

mod_update:
	GO111MODULE=on go mod download

mod_init:
	go mod init src

get_bin:
	GO111MODULE=on go get github.com/pilu/fresh

build:
	go build -o bin/server main.go

clean:
	rm -rf bin/*

fmt:
	@find src -type f -iregex '.*\.go' | xargs -L 1 -P 4 gofmt

imports:
	find src -type f -iregex '.*\.go' | xargs -L 1 -P 2 goimports -w
	find tests -type f -iregex '.*\.go' | xargs -L 1 -P 2 goimports -w

serve:
	#@sysctl -w fs.inotify.max_user_instances=65535
	@fresh -c ./fresh.conf

.PHONY: mod_update mod_init build fmt imports clean serve 
