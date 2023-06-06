export GOPROXY=https://goproxy.io
export GO111MODULE=on

HOMEDIR := $(shell pwd)

all: mod build

mod:
	go mod tidy -v

build:
	bash $(HOMEDIR)/build.sh

config:
	ibax-cli config
console:
	ibax-cli console

startup: config console
