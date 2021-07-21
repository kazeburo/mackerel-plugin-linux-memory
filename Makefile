VERSION=0.0.2
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION}"
all: mackerel-plugin-linux-memory

.PHONY: mackerel-plugin-linux-process-status

mackerel-plugin-linux-memory: main.go
	go build $(LDFLAGS) -o mackerel-plugin-linux-memory

linux: main.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-linux-memory

fmt:
	go fmt ./...

check:
	go test ./...

clean:
	rm -rf mackerel-plugin-linux-memory

tag:
	git tag v${VERSION}
	git push origin v${VERSION}
	git push origin main
