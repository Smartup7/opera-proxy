PROGNAME = opera-proxy-checker
OUTSUFFIX = bin/$(PROGNAME)
BUILDOPTS = -a -tags netgo -trimpath -asmflags -trimpath
LDFLAGS = -ldflags '-s -w -extldflags "-static"'
LDFLAGS_NATIVE = -ldflags '-s -w'

src = $(wildcard *.go)

all: bin-native bin-linux-arm64

bin-native: $(OUTSUFFIX)
bin-linux-arm64: $(OUTSUFFIX).linux-arm64

$(OUTSUFFIX): $(src)
	go build $(LDFLAGS_NATIVE) -o $@

$(OUTSUFFIX).linux-arm64: $(src)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build $(BUILDOPTS) $(LDFLAGS) -o $@

clean:
	rm -f bin/*

fmt:
	go fmt ./...

install:
	go install $(LDFLAGS_NATIVE) .

.PHONY: clean all native fmt install \
	bin-native \
	bin-linux-arm64