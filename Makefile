# Define service names
SERVICES := shortly-api-service shortly-kgs-service

# Define build directory (relative to root)
BUILD_DIR := bin

.PHONY: build run migrate clean help

help:
	@echo "Usage: make [command] SERVICE=<service_name>"
	@echo "Commands:"
	@echo "  build    Build the specified service"
	@echo "  run      Build and run the specified service"
	@echo "  migrate  Run migrations for the specified service"
	@echo "  clean    Remove built binaries"
	@echo ""
	@echo "Available Services: $(SERVICES)"

# Build a specific service by changing into its directory first
build:
	@if [ -z "$(SERVICE)" ]; then \
		echo "❌ SERVICE variable is required. Run: make build SERVICE=<service_name>"; exit 1; \
	fi
	mkdir -p $(BUILD_DIR)
	rm -rf $(BUILD_DIR)/$(SERVICE)  
	# Change directory into the service folder, then build
	cd services/$(SERVICE) && \
	go build -o ../../$(BUILD_DIR)/$(SERVICE) ./cmd/main.go

# Run a specific service (builds first)
run: build
	@if [ -z "$(SERVICE)" ]; then \
		echo "❌ SERVICE variable is required. Run: make run SERVICE=<service_name>"; exit 1; \
	fi
	# Change to service directory before running
	cd services/$(SERVICE) && ../../$(BUILD_DIR)/$(SERVICE)


# Run migrations for a specific service by changing into its directory first
migrate:
	@if [ -z "$(SERVICE)" ]; then \
		echo "❌ SERVICE variable is required. Run: make migrate SERVICE=<service_name>"; exit 1; \
	fi
	cd services/$(SERVICE) && \
	go run internal/migrations/migration.go

clean:
	rm -rf $(BUILD_DIR)
