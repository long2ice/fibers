PROJECT_NAME := "fibers"
PKG := "github.com/long2ice/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/ )

lint:
	@golint -set_exit_status ${PKG_LIST}

dep:
	@go get -v -d ./...

build: dep
	@go build -v $(PKG)

clean:
	@rm -f $(PROJECT_NAME)

format:
	@gofumpt -l -w . && golines . -w