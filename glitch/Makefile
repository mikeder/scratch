
wasm:
#   cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
	GOOS=js GOARCH=wasm go build -ldflags "-s" -o ./bin/gophermark.wasm cmd/gophermark/main.go

compress:
	ls -lh ./bin/gophermark.wasm
	gzip -f --best -c ./bin/gophermark.wasm > ./bin/gophermark.wasm.gz
	ls -lh ./bin/gophermark.wasm

serve:
	cp ./bin/gophermark.wasm ./cmd/server/assets/
	go run cmd/server/server.go

run: wasm compress serve
