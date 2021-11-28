.PHONY: container clean

container:
	docker build -t domprinter ./

build/%: %/main.go
	CGO_ENABLED=0 GOOS=linux go build -o $@ $<

clean:
	rm -rf build
