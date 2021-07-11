NAME=satori
BINARY=terraform-provider-${NAME}

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
	mkdir -p ~/.terraform.d/plugins
	mv ${BINARY} ~/.terraform.d/plugins/

.PHONY: test
test:
	go test -i $(TEST) || exit 1
	go test ./... $(TESTARGS) -timeout=30s -parallel=4

# Run acceptance tests
.PHONY: testacc
testacc: build
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
