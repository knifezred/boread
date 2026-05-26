.PHONY: run build build-local test clean lint tidy swag

APP_NAME    := boread
BUILD_DIR   := bin
CMD_DIR     := ./cmd/api

GOOS        := linux
GOARCH      := amd64
CGO_ENABLED ?= 0

LDFLAGS     := -s -w

run:
	go run $(CMD_DIR)/

build: tidy
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
		go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/
	@echo "Built: $(BUILD_DIR)/$(APP_NAME) ($(GOOS)/$(GOARCH))"

build-local:
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/
	@echo "Built: $(BUILD_DIR)/$(APP_NAME) (local)"

test:
	go test -v -count=1 ./...

lint:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -rf $(BUILD_DIR)

swag:
	swag init -g $(CMD_DIR)/main.go -o ./docs --parseDependency --parseInternal