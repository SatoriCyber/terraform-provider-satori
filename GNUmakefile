NAME=satori
HOSTNAME=satoricyber.com
NAMESPACE=terraform
VERSION=1.0.0
BINARY=terraform-provider-${NAME}
OS_ARCH:=$(shell uname -s)-$(shell uname -m)
LOCAL_INSTALL_DIR=~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

default: testacc

.PHONY: init
init:
	go mod vendor
	go mod tidy

.PHONY: docs
docs:
	go generate

.PHONY: build
build:
	go build -o ${BINARY}

.PHONY: install
install: build
	mkdir -p ${LOCAL_INSTALL_DIR}
	mv ${BINARY} ${LOCAL_INSTALL_DIR}

.PHONY: test
test:
	go test -i $(TEST) || exit 1
	go test ./... $(TESTARGS) -timeout=30s -parallel=4

# Run acceptance tests
.PHONY: testacc
testacc: build
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
