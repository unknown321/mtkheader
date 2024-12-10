PRODUCT=mtkheader
GOOS=linux
GOARCH=$(shell uname -m)
GOARM=
NAME=$(PRODUCT)-$(GOOS)-$(GOARCH)$(GOARM)

ifeq ($(GOARCH),x86_64)
	override GOARCH=amd64
endif

$(NAME):
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a \
		-ldflags "-w -s" \
		-trimpath \
		-o $(NAME)

build: $(NAME)

test:
	go test -v ./...

all: test
	$(MAKE) build
	$(MAKE) GOARCH=arm GOARM=5 build

clean:
	rm -rfv $(PRODUCT)-*

release:
	$(MAKE) build
	$(MAKE) GOARCH=arm GOARM=5 GOOS=linux build

.PHONY: test
.DEFAULT_GOAL := all
