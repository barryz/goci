NOW=$(shell date '+%Y-%m-%d_%H:%M:%S')
REV?=$(shell git rev-parse HEAD)
LDFLAGS=-ldflags '-X main.Build=${NOW}@${REV}'
BINARY="goci"

ifeq ($(GO111MODULE),on)
	OPTS=-mod=vendor
else
	OPTS=
endif

default: build

build:
	go build ${OPTS} -o ./${BINARY} ${LDFLAGS}

test:
	go test ${OPTS} ./...

clean:
	if test -f ${BINARY}; then \
	rm -f ${BINARY}; fi


.PHONY: build clean
