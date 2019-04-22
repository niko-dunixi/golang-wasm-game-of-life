.PHONY: build-client build-server build clean gomodgen run

build:	clean build-client build-server

build-client: gomodgen
	export GO111MODULE=on
	env GOARCH=wasm GOOS=js go build -ldflags="-s -w" -o bin/client.wasm client/main.go

build-server: build-client
	export GO111MODULE=on
	go build -ldflags="-s -w" -o bin/server server/main.go
	cp server/index.html bin/index.html
	cp server/wasm_exec.js bin/wasm_exec.js

run:	build

clean:
	rm -rf ./bin ./vendor Gopkg.lock

go.mod:
	chmod u+x gomod.sh
	./gomod.sh
