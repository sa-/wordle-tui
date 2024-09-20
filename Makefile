.PHONY: demo
demo:
	vhs demo/demo.tape -o demo/demo.gif

.PHONY: run
run:
	go run main.go
