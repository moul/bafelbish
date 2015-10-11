.PHONY: convey
convey:
	go get github.com/smartystreets/goconvey
	goconvey -cover -port=9045 -workDir="$(realpath .)" -depth=0
