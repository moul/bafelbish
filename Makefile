all: build


.PHONY: build
build:
	go build -o bafelbish ./cmd/bafelbish


.PHONY: convey
convey:
	go get github.com/smartystreets/goconvey
	goconvey -cover -port=9045 -workDir="$(realpath .)" -depth=0
