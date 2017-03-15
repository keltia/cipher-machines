# Makefile

VPATH=  config:lib:cmd/cm-test
GOBIN=  ${GOPATH}/bin
OPTS=   -ldflags="-s -w" -v
SRCS=   rotor.go enigma.go sigaba.go main.go
TESTS=	rotor_test.go enigma_test.go

all:	cm-test

cm-test:	$(SRCS)
	go build ${OPTS} ./cmd/...
	go test -v .

tests:	$(TESTS)
	go test -v .
