NOW=$(shell date '+%Y-%m-%d_%H:%M:%S')
REV?=$(shell git rev-parse HEAD)
LDFLAGS=-ldflags '-X main.Build=${NOW}@${REV}'
BINARY="goci"

build:
	go build -o ./${BINARY} ${LDFLAGS}

clean:
	if test -f ${BINARY}; then \
	rm -f ${BINARY}; fi

.PHONY: build clean
