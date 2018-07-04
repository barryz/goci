NOW=$(shell date '+%Y-%m-%d_%H:%M:%S')
REV?=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags '-X github.com/barryz/goci/main.Build=${NOW}@${REV} -w -s'
BINARY="goci"

build:
	go build -o ./${BINARY} ${LDFLAGS}

clean:
	if test -f ${BINARY}; then \
	rm -f ${BINARY}; fi

.PHONY: build clean
