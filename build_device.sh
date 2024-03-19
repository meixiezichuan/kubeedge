 CGO_ENABLED=1  go build -v -o _output/local/bin -ldflags="${GO_LDFLAGS} -w -s -extldflags -static" ./edgedevice/cmd/edgedevice
