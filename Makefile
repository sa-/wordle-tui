.PHONY: demo
demo:
	go build -o out/wordle
	PATH="$$(pwd)/out:$$PATH" vhs demo/demo.tape -o demo/demo.gif

.PHONY: build
build:
	go build -o out/wordle

.PHONY: run
run: build
	./out/wordle

.PHONY: clean
clean:
	rm -r out
