GFLAGS=-v -trimpath
BUILD=go build

BINARY=pomogo

all: $(BINARY)

$(BINARY): 
	go mod tidy
	go install

clean: 
	rm -rf ./$(BINARY)

build:
	$(BUILD) $(GFLAGS) -o $@
