BIN := elect

.PHONY: all
all: build

.PHONY: build
build:
	go build -i -o ./build/${BIN} ./elect

.PHONY: clean
clean:
	rm -rf ./build

.PHONY: install
install:
	go install -v ./elect
