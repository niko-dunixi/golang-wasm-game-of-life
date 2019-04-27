.PHONY: build-client build-server build clean gomodgen run

build:	clean build-client build-server

build-client: gomodgen
	export GO111MODULE=on
	env GOARCH=wasm GOOS=js go build -ldflags="-s -w" -o bin/client.wasm .

build-server: build-client
	export GO111MODULE=on
	go build -ldflags="-s -w" -o bin/server server/main.go
	cp server/index.html bin/index.html
	cp server/wasm_exec.js bin/wasm_exec.js

run:	build
	./bin/server

clean:
	rm -rf ./bin ./vendor Gopkg.lock

go.mod:
	go mod init "github.com/paul-nelson-baker/wasm-game-of-life"
