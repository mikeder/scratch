.DEFAULT_GOAL = wasm

build_targets := gophermark game

wasm_exec.js:
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" ./cmd/server/assets/scripts/

wasm: wasm_exec.js
	$(foreach target,$(build_targets), \
		GOOS=js GOARCH=wasm go build -ldflags "-s" -o ./bin/$(target).wasm cmd/$(target)/main.go; \
		cp ./bin/$(target).wasm ./cmd/server/assets/; \
	)

compress:
	ls -lh ./bin/gophermark.wasm
	gzip -f --best -c ./bin/gophermark.wasm > ./bin/gophermark.wasm.gz
	ls -lh ./bin/gophermark.wasm

serve:
	go run cmd/server/server.go

run: wasm compress serve
