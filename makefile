all:

	make mac
	make linux

mac:
	GOOS="darwin" GOARCH=amd64
	go build hollywood.go

linux:
	export GOARCH="amd64"
	export GOOS="linux"
	go build hollywood.go