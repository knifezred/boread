.PHONY: run-backend run-frontend build build-local build-all build-frontend build-backend pack test clean lint tidy swag seed help

# ========== Project Paths ==========
PROJECT_ROOT := $(CURDIR)/app
BACKEND_DIR  := $(PROJECT_ROOT)/server
FRONTEND_DIR := $(PROJECT_ROOT)/web

# ========== Output Paths ==========
BUILD_DIR        := $(PROJECT_ROOT)/build
BACKEND_OUTPUT   := $(BUILD_DIR)/rootfs_amd64/bin/boread
BACKEND_LOCAL    := $(BUILD_DIR)/rootfs_amd64/bin/boread.exe
FRONTEND_OUTPUT  := $(BUILD_DIR)/rootfs_common/www

# ========== Go Build Flags ==========
GOOS        := linux
GOARCH      := amd64
CGO_ENABLED ?= 0
LDFLAGS     := -s -w
BUILD_NUM   ?= 1

# ========== Windows Commands ==========
RM := powershell -Command "Remove-Item -Recurse -Force -ErrorAction SilentlyContinue"
MKDIR := powershell -Command "New-Item -ItemType Directory -Force -Path"

# ========== Development ==========
run-backend:
	@echo "[INFO] Starting backend server..."
	@echo "[INFO] Backend: http://localhost:8080"
	cd $(BACKEND_DIR) && go run ./cmd/api/

run-frontend:
	@echo "[INFO] Starting frontend server..."
	@echo "[INFO] Frontend: http://localhost:5173"
	cd $(FRONTEND_DIR) && pnpm run dev

# ========== Backend Build ==========
build-backend: tidy
	@echo "[INFO] Building Go backend ($(GOOS)/$(GOARCH))..."
	$(MKDIR) $(BUILD_DIR)/rootfs_amd64/bin > /dev/null 2>&1
	cd $(BACKEND_DIR) && export CGO_ENABLED=$(CGO_ENABLED) && export GOOS=$(GOOS) && export GOARCH=$(GOARCH) && \
		go build -ldflags "$(LDFLAGS)" -o $(BACKEND_OUTPUT) ./cmd/api/
	@echo "[OK] Backend built: $(BACKEND_OUTPUT)"

build-local:
	@echo "[INFO] Building Go backend (local)..."
	$(MKDIR) $(BUILD_DIR)/rootfs_amd64/bin > /dev/null 2>&1
	cd $(BACKEND_DIR) && go build -ldflags "$(LDFLAGS)" -o $(BACKEND_LOCAL) ./cmd/api/
	@echo "[OK] Backend built: $(BACKEND_LOCAL)"

tidy:
	@echo "[INFO] Running go mod tidy..."
	cd $(BACKEND_DIR) && go mod tidy

# ========== Frontend Build ==========
build-frontend:
	@echo "[INFO] Building Vue frontend..."
	$(MKDIR) $(FRONTEND_OUTPUT) > /dev/null 2>&1
	cd $(FRONTEND_DIR) && pnpm install
	cd $(FRONTEND_DIR) && pnpm run build
	@echo "[OK] Frontend built: $(FRONTEND_OUTPUT)"

build-frontend-clean:
	@echo "[INFO] Cleaning frontend output..."
	$(RM) $(FRONTEND_OUTPUT)
	$(MKDIR) $(FRONTEND_OUTPUT) > /dev/null 2>&1
	@echo "[INFO] Building Vue frontend..."
# 	cd $(FRONTEND_DIR) && pnpm install
	cd $(FRONTEND_DIR) && pnpm run build
	@echo "[OK] Frontend built: $(FRONTEND_OUTPUT)"

# ========== Pack ==========
pack:
	@echo "[INFO] Running ugcli pack..."
	cd $(BUILD_DIR) && ugcli pack --arch amd64 --build $(BUILD_NUM)
	@echo "[OK] Pack completed"

# ========== One-Command Build ==========
build-all: build-frontend-clean build-backend
	@echo "[OK] Complete build finished!"
	@echo "[OK] Backend:  $(BACKEND_OUTPUT)"
	@echo "[OK] Frontend: $(FRONTEND_OUTPUT)"

build-all-local: build-frontend-clean build-local
	@echo "[OK] Local build finished!"
	@echo "[OK] Backend:  $(BACKEND_LOCAL)"
	@echo "[OK] Frontend: $(FRONTEND_OUTPUT)"

build: build-backend

# ========== Testing & Quality ==========
test:
	cd $(BACKEND_DIR) && go test -v -count=1 ./...

lint:
	cd $(BACKEND_DIR) && go vet ./...

# ========== Docs & Data ==========
swag:
	cd $(BACKEND_DIR) && swag init -g ./cmd/api/main.go -o ./docs --parseDependency --parseInternal

seed:
	cd $(BACKEND_DIR) && go run ./cmd/api/ -seed

# ========== Clean ==========
clean: clean-backend clean-frontend
	@echo "[OK] All cleaned"

clean-backend:
	@echo "[INFO] Cleaning backend..."
	$(RM) $(BUILD_DIR)/rootfs_amd64
	@echo "[OK] Backend cleaned"

clean-frontend:
	@echo "[INFO] Cleaning frontend..."
	$(RM) $(FRONTEND_OUTPUT)
	@echo "[OK] Frontend cleaned"

# ========== Help ==========
help:
	@echo "Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  make run-backend  - Start backend only (http://localhost:8080)"
	@echo "  make run-frontend - Start frontend only (http://localhost:5173)"
	@echo ""
	@echo "Build:"
	@echo "  make build-all        - Build frontend + backend (Linux)"
	@echo "  make build-all-local  - Build frontend + backend (Windows)"
	@echo "  make build-frontend   - Build frontend only"
	@echo "  make build-backend    - Build backend only"
	@echo "  make pack BUILD_NUM=1  - Pack with build version 1"
	@echo ""
	@echo "Clean:"
	@echo "  make clean - Clean all build outputs"
	@echo ""
	@echo "Other:"
	@echo "  make test  - Run tests"
	@echo "  make swag  - Generate Swagger docs"
	@echo "  make seed  - Seed database"