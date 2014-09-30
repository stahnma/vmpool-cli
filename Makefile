

build: vmpool

fmt:
	go fmt vmpool.go

vmpool:
	go build -ldflags "-X main.version `date -u +%Y%m%d%.H%M%S`" vmpool.go


install:
	sudo cp -pr vmpool $(DESTDIR)/usr/local/bin

linux:
# In order to get the cross-compile options on mac, install via
#     brew install go --cross-compile-common
	GOARCH=amd64 GOOS=linux go build

clean:
	rm -rf vmpool

.PHONY: intall fmt clean
