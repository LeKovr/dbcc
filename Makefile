
SHELL    = /bin/bash
SRC_LIST := main.go,dbcc.go
PRG      ?= $(shell basename $$PWD)
LOG      ?= $(PRG).log
PID      ?= /tmp/$(PRG).pid
HOST     ?= localhost
PORT     ?= 8000
ARCH     ?= amd64
OSARCH   ?= linux/$(ARCH)
ALLARCH  ?= "linux/amd64 linux/386 windows/amd64 darwin/386"
GOPATH   ?=
DBUSER   ?= op
APPUSER  ?= appuser
APPPASS  ?= apppass
KEYFILE  ?= .appkey

all: chkgo clean build test sum

dist: chkgo clean buildall sum

chkgo:
	@echo "*** $@ ***"
	@[ -f $(GOPATH)/bin/gox ] || { echo "Error: gox compiler not found" ; exit 1 ; }

clean:
	@echo "*** $@ ***"
	@for f in bin/$(PRG)*; do [ -e "$$f" ] && rm "$$f" || echo "no bin files" ; done
	@for f in $(PRG)*.zip; do [ -e "$$f" ] && rm "$$f" || echo "no zip files" ; done
	@[ -d bin ] && [ -f bin/SHA256SUMS ] && rm bin/SHA256SUMS || true 
	@[ -d bin ] && rm -d bin || true

build:
	@echo "*** $@ ***"
	@[ -d bin ] || mkdir bin
	@pushd bin ; \
	gox -osarch="$(OSARCH)" ../ && popd || { popd ; exit 1 ; }

buildall:
	@echo "*** $@ ***"
	@[ -d bin ] || mkdir bin
	@pushd bin ; \
	gox -osarch=$(ALLARCH) ../ && popd || { popd ; exit 1 ; }

buildpack:
	@echo "*** $@ ***"
	@pushd bin ; \
	for f in $(PRG)* ; do [ -e "$$f" ] && zip "../$$f.zip" "$$f" || echo "no files" ; done

test:
	@echo "*** $@ ***"
	go test

testall:
	@echo "*** $@ ***"
	@PGUSER=$(DBUSER) DBCC_TEST_DB=1 go test

sum:
	@echo "*** $@ ***"
	@[ -d bin ] || mkdir bin
	@pushd bin ; sha256sum $(PRG)* > SHA256SUMS ; popd

get-deps:
	go get github.com/mitchellh/gox
	go get -t

appkey:
	@echo "*** $@ ***"
	@[ -e $(KEYFILE) ] || { LC_ALL=C < /dev/urandom tr -dc _A-Z-a-z-0-9 | head -c$${1:-16} > $(KEYFILE) ; echo "$(KEYFILE) created." ; }

restart: stop start

start:
	@echo "*** $@ ***"
	@nohup bin/$(PRG)_linux_$(ARCH) -host $(HOST) -port $(PORT) >>$(LOG) 2>&1 & echo $$! > $(PID)
	@echo "Started, pid=`cat $(PID)`"

stop:
	@echo "*** $@ ***"
	@[ -f $(PID) ] && kill `cat $(PID)` || echo "No pidfile"
	@[ -f $(PID) ] && rm $(PID) || true

status:
	@echo "*** $@ ***"
	@[ -f $(PID) ] && kill -0 `cat $(PID)` && echo "running" || echo "No such process"

run: appkey
	@echo "*** $@ ***"
	@PGUSER=$(DBUSER) APP_KEY=$$(cat $(KEYFILE)) go run {$(SRC_LIST)} -host $(HOST) -port $(PORT)

use:
	@echo "*** $@ ***"
	curl "http://$(HOST):$(PORT)/?key=$$(cat $(KEYFILE))&name=$(APPUSER)&pass=$(APPPASS)"

.PHONY: all stop start chkgo rm build buildall sum run
