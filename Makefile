# NOTE(timonwong): From now on, this Makefile only works on go1.10+
GO    := go

REPO_PATH               ?= github.com/imperfectgo/common
TESTARGS                ?= -v -race
COVERARGS               ?= -covermode=atomic
TEST                    ?= $(shell go list ./... | grep -v '/vendor/')
TESTPKGS                ?= $(shell go list ./... | grep -v '/cmd/')
GOFMT_FILES             ?= $(shell find . -name '*.go' | grep -v vendor | xargs)
FIRST_GOPATH            := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))
DEP                     := $(FIRST_GOPATH)/bin/dep
OVERALLS                := $(FIRST_GOPATH)/bin/overalls
GOIMPORTS               := $(FIRST_GOPATH)/bin/goimports
GOMETALINTER            := $(FIRST_GOPATH)/bin/gometalinter

export REPO_PATH

_comma := ,
_space :=
_space +=


.PHONY: all
all: format test


$(DEP):
	@echo ">> installing golang dep tool"
	@$(GO) get -u "github.com/golang/dep/cmd/dep"


$(OVERALLS):
	@echo ">> installing overalls tool"
	@$(GO) get -u "github.com/go-playground/overalls"


$(GOIMPORTS):
	@echo ">> installing goimports tool"
	@$(GO) get -u "go get golang.org/x/tools/cmd/goimports"


$(GOMETALINTER):
	@echo ">> installing gometalinter"
	@$(GO) get -u "github.com/alecthomas/gometalinter"
	@$(GOMETALINTER) --install --update


.PHONY: dep
dep: $(DEP)
	@dep ensure


.PHONY: test
test:
	@echo ">> running tests"
	@$(GO) test $(TEST) $(TESTARGS)


.PHONY: cover
cover: $(OVERALLS)
	@echo ">> running test coverage"
	@rm -f coverage.txt
	@$(OVERALLS) -project=$(REPO_PATH) $(COVERARGS) -- -coverpkg=./... $(TESTARGS) && \
		mv overalls.coverprofile coverage.txt


.PHONY: lint
lint: $(GOMETALINTER)
	@echo ">> linting code"
	@$(GOMETALINTER) --vendor --disable-all \
		--enable=varcheck \
		--enable=gosimple \
		--enable=misspell \
		--enable=vet \
		--enable=vetshadow \
		--enable=golint \
		--deadline=10m \
		./...


.PHONY: fmtcheck
fmtcheck: $(GOMETALINTER)
	@echo ">> checking code style"
	@$(GOMETALINTER) --vendor --disable-all \
		--enable=gofmt \
		--enable=goimports \
		./...


.PHONY: format
format: $(GOIMPORTS)
	@echo ">> formatting code"
	@$(GOIMPORTS) -local "$(REPO_PATH)" -w $(GOFMT_FILES)
