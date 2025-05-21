all: build

build:
	GOARCH=arm64 go build -o smalltalk-vm

run: build
	./smalltalk-vm

clean:
	rm -f smalltalk-vm

.PHONY: all build run clean
