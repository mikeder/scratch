.DEFAULT_GOAL = copy

build_targets := ebitengine/jam2024 ebitengine/kage glitch

wasm_exec.js:
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" ./server/assets/shared/

wasm: wasm_exec.js
	$(foreach target,$(build_targets), \
		GOOS=js GOARCH=wasm go build -ldflags "-s" -o ./bin/$(target).wasm $(target)/main.go; \
	)

compress:
	$(foreach target,$(build_targets), \
		ls -lh ./bin/$(target).wasm; \
		gzip -f --best -c ./bin/$(target).wasm > ./bin/$(target).wasm.gz; \
		ls -lh ./bin/$(target).wasm; \
	)

copy: wasm
	$(foreach target,$(build_targets), \
		cp ./bin/$(target).wasm ./server/assets/$(target)/; \
	)

serve:
	go run server/server.go

run: wasm copy serve

publish:
	sh scripts/publish_to_ghpages.sh