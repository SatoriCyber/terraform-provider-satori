NAME=satori
HOSTNAME=satoricyber.com
NAMESPACE=terraform
VERSION=1.0.0
BINARY=terraform-provider-${NAME}
# This piece of code defines the OS for the machines that installs the provider.
UNAME_P := $(shell uname -p)
ifeq ($(filter %86,$(UNAME_P)),)
  OS_ARCH=darwin_amd64
endif
ifneq ($(filter arm%,$(UNAME_P)),)
  OS_ARCH=darwin_arm64
endif
$(info Current architecture (OS_ARCH) is $(OS_ARCH))
LOCAL_INSTALL_DIR=~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
$(info Installation folder (LOCAL_INSTALL_DIR) is $(LOCAL_INSTALL_DIR))

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

.PHONY: vuln
vuln:
	go tool govulncheck