GFLAGS=-v -trimpath
BUILD=go build

BINARY=pomogo

all: $(BINARY)

$(BINARY): 
	go mod tidy
	$(BUILD) $(GFLAGS) -o $@

clean: 
	rm -rf ./$(BINARY)
