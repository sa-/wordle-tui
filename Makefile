.PHONY: demo
demo: build
	PATH="$$(pwd)/out:$$PATH" WORDLE_WORD=charm vhs demo/demo.tape -o demo/demo.gif

.PHONY: build
build:
	go build -ldflags "-s -w" -o out/wordle-tui

.PHONY: run
run: build
	./out/wordle-tui

.PHONY: clean
clean:
	rm -r out
