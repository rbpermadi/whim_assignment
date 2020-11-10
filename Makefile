TIMESTAMP = $(shell date +'%Y%m%d%H%M%S')
PWD	      = $(shell pwd)

IMAGE_NAME = rbpermadi/whim-assignment
VERSION    = $(shell git show -q --format=%h)

# LOCAL SETUP FOR COMPILED BINARY
GOOS	 ?= linux
GOARCH  = amd64
ODIR    = _output

export GO111MODULE ?= on

run:bin
	./_output/whim

bin:
	go build -o _output/whim app/web-service/main.go

dep:
	go mod tidy

test:
	go test -v ./...

coverage:
	go test -coverprofile fmt ./...

compile:
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o $(ODIR)/whim_docker app/web-service/main.go

build:
	docker build -t $(IMAGE_NAME) -f ./Dockerfile .