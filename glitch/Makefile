.DEFAULT_GOAL = wasm

wasm_exec.js:
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" ./cmd/server/assets/scripts/

wasm: wasm_exec.js
	GOOS=js GOARCH=wasm go build -ldflags "-s" -o ./bin/game.wasm cmd/game/main.go
	cp ./bin/game.wasm ./cmd/server/assets/
	
	GOOS=js GOARCH=wasm go build -ldflags "-s" -o ./bin/gophermark.wasm cmd/gophermark/main.go
	cp ./bin/gophermark.wasm ./cmd/server/assets/

compress:
	ls -lh ./bin/gophermark.wasm
	gzip -f --best -c ./bin/gophermark.wasm > ./bin/gophermark.wasm.gz
	ls -lh ./bin/gophermark.wasm

serve:
	go run cmd/server/server.go

run: wasm compress serve
