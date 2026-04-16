BINARY := bin/TimeKeeper
INSTALL_PATH := /usr/local/bin/TimeKeeper
SRC := ./src/

.PHONY: all build clean install

all: build

build:
	go build -o $(BINARY) $(SRC)

clean:
	rm -f $(BINARY)

install: build
	cp $(BINARY) $(INSTALL_PATH)
