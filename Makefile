PRODUCT=mtkheader
GOOS=linux
GOARCH=$(shell uname -m)
GOARM=
NAME=$(PRODUCT)-$(GOOS)-$(GOARCH)$(GOARM)

ifeq ($(GOARCH),x86_64)
	override GOARCH=amd64
endif

build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a \
		-ldflags "-w -s" \
		-trimpath \
		-o $(NAME)

build-arm: GOARCH=arm
build-arm: GOARM=5
build-arm: build

test:
	go test -v ./...

all: test
	$(MAKE) build
	$(MAKE) build-arm

clean:
	rm -rfv $(PRODUCT)-*

.PHONY: test
.DEFAULT_GOAL := all
