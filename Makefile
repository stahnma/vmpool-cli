

build: vmpool

fmt:
	go fmt vmpool.go

vmpool:
	go build vmpool.go


install:
	sudo cp -pr vmpool /usr/local/bin


clean:
	rm -rf vmpool
