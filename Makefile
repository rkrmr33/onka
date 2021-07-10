VERSION= $(shell cat VERSION)

SERVER_BIN_NAME := onka-server
DAEMON_BIN_NAME := onkad
CLI_BIN_NAME := onka

OUT_DIR := dist
PROTO_GEN=gofast

ifdef BIN_NAME
SRCS := $(shell echo cmd/$(BIN_NAME)/*.go && go list -f '{{ join .Deps "/*.go\n" }}' ./cmd/$(BIN_NAME) | grep 'rkrmr33/onka' | cut -c 25-)
endif

ifndef GOBIN
ifndef GOPATH
$(error GOPATH is not set, please make sure you set your GOPATH correctly!)
endif
GOBIN=$(GOPATH)/bin
ifndef GOBIN
$(error GOBIN is not set, please make sure you set your GOBIN correctly!)
endif
endif

def:
	@make local BIN_NAME=$(DAEMON_BIN_NAME)
	@make local BIN_NAME=$(SERVER_BIN_NAME)
	@make local BIN_NAME=$(CLI_BIN_NAME)

### Daemon ###
PHONY: $(DAEMON_BIN_NAME)
$(DAEMON_BIN_NAME):
	@make all BIN_NAME=$(DAEMON_BIN_NAME)

PHONY: $(DAEMON_BIN_NAME)-local
$(DAEMON_BIN_NAME)-local:
	@make local BIN_NAME=$(DAEMON_BIN_NAME)

### Server ###
PHONY: $(SERVER_BIN_NAME)
$(SERVER_BIN_NAME):
	@make all BIN_NAME=$(SERVER_BIN_NAME)

PHONY: $(SERVER_BIN_NAME)-local
$(SERVER_BIN_NAME)-local:
	@make local BIN_NAME=$(SERVER_BIN_NAME)

### CLI ###
PHONY: $(CLI_BIN_NAME)
$(CLI_BIN_NAME):
	@make all BIN_NAME=$(CLI_BIN_NAME)

PHONY: $(CLI_BIN_NAME)-local
$(CLI_BIN_NAME)-local:
	@make local BIN_NAME=$(CLI_BIN_NAME)


### Generic ###
PHONY: all
all: $(OUT_DIR)/$(BIN_NAME)-linux-amd64 $(OUT_DIR)/$(BIN_NAME)-darwin-amd64 $(OUT_DIR)/$(BIN_NAME)-windows-amd64 $(OUT_DIR)/$(BIN_NAME)-linux-arm64 $(OUT_DIR)/$(BIN_NAME)-linux-ppc64le $(OUT_DIR)/$(BIN_NAME)-linux-s390x

PHONY: local
local: $(OUT_DIR)/$(BIN_NAME)-$(shell go env GOOS)-$(shell go env GOARCH)
	@cp $(OUT_DIR)/$(BIN_NAME)-$(shell go env GOOS)-$(shell go env GOARCH) /usr/local/bin/$(BIN_NAME)
	
$(OUT_DIR)/$(BIN_NAME)-linux-amd64: GO_FLAGS='GOOS=linux GOARCH=amd64 CGO_ENABLED=0'
$(OUT_DIR)/$(BIN_NAME)-darwin-amd64: GO_FLAGS='GOOS=darwin GOARCH=amd64 CGO_ENABLED=0'
$(OUT_DIR)/$(BIN_NAME)-windows-amd64: GO_FLAGS='GOOS=windows GOARCH=amd64 CGO_ENABLED=0'
$(OUT_DIR)/$(BIN_NAME)-linux-arm64: GO_FLAGS='GOOS=linux GOARCH=arm64 CGO_ENABLED=0'
$(OUT_DIR)/$(BIN_NAME)-linux-ppc64le: GO_FLAGS='GOOS=linux GOARCH=ppc64le CGO_ENABLED=0'
$(OUT_DIR)/$(BIN_NAME)-linux-s390x: GO_FLAGS='GOOS=linux GOARCH=s390x CGO_ENABLED=0'

$(OUT_DIR)/$(BIN_NAME)-%: $(SRCS)
	@GO_FLAGS=$(GO_FLAGS) \
	BINARY_NAME=$(BIN_NAME) \
	VERSION=$(VERSION) \
	OUT_FILE=$(OUT_DIR)/$(BIN_NAME)-$* \
	MAIN=./cmd/$(BIN_NAME) \
	./hack/build.sh

PHONY: gen
gen: $(GOBIN)/mockery genproto gengrammer
	@echo generating mocks...
	@go generate ./...

PHONY: genproto
genproto: /usr/local/bin/protoc $(OUT_DIR)/protoc-gen-$(PROTO_GEN) $(OUT_DIR)/protoc-gen-go-grpc
	@echo generating protobuf...
	@PROTO_GEN=$(PROTO_GEN) ./hack/generate.sh
	@go mod tidy

PHONY: gengrammer
gengrammer: $(GOBIN)/pigeon $(GOBIN)/goimports
	@echo generating grammer...
	@go generate ./jag/compiler/bexpr

.PHONY: release
release: clean-worktree
	@./hack/release.sh

.PHONY: clean-worktree
clean-worktree:
	@./hack/clean-worktree.sh

PHONY: lint
lint: $(GOBIN)/golangci-lint
	@echo linting go code...
	@GOGC=off golangci-lint run --fix --timeout 6m
	
.PHONY: test
test:
	./hack/test.sh

PHONY: clean
clean: 
	@rm -rf $(OUT_DIR) vendor || true

PHONY: clean-pb
clean-pb:
	@find . -type f -name '*.pb.go' -exec rm {} +

./vendor:
	@echo vendoring...
	@go mod vendor

$(GOBIN)/mockery:
	@mkdir dist || true
	@echo installing: mockery
	@curl -L -o dist/mockery.tar.gz -- https://github.com/vektra/mockery/releases/download/v1.1.1/mockery_1.1.1_$(shell uname -s)_$(shell uname -m).tar.gz
	@tar zxvf dist/mockery.tar.gz mockery
	@chmod +x mockery
	@mkdir -p $(GOBIN)
	@mv mockery $(GOBIN)/mockery
	@mockery -version

$(GOBIN)/golangci-lint:
	@mkdir dist || true
	@echo installing: golangci-lint
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.36.0

/usr/local/bin/protoc:
	@echo downloading protoc...
	@curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/$(PROTOC_PKG)
	@sudo unzip -o $(PROTOC_PKG) -d /usr/local bin/protoc
	@sudo unzip -o $(PROTOC_PKG) -d /usr/local 'include/*'
	@rm -f $(PROTOC_PKG)

$(OUT_DIR)/protoc-gen-$(PROTO_GEN): ./vendor
	@go build -o dist/protoc-gen-$(PROTO_GEN) ./vendor/github.com/gogo/protobuf/protoc-gen-$(PROTO_GEN)

$(OUT_DIR)/protoc-gen-go-grpc: ./vendor
	@go build -o dist/protoc-gen-go-grpc ./vendor/google.golang.org/grpc/cmd/protoc-gen-go-grpc

$(GOBIN)/pigeon: CUR=$(shell pwd)
$(GOBIN)/pigeon:
	@cd $(shell mktemp -d)
	@echo downloading pigeon...
	@go get github.com/mna/pigeon
	@cd $(CUR)

$(GOBIN)/goimports: CUR=$(shell pwd)
$(GOBIN)/goimports:
	@cd $(shell mktemp -d)
	@echo downloading goimports...
	@go get golang.org/x/tools/cmd/goimports
	@cd $(CUR)
