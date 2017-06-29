GO=$(shell which go)

build:
	# CGO_ENABLED=0 makes it possible to use the binary with libc libraries that
	# aren't the libc in use on the system used to build the binary.
	CGO_ENABLED=0 $(GO) build wkhtmltox-server.go

run: build
	$(PWD)/wkhtmltox-server
