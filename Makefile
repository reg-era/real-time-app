# Application-related variables
PORT           = 8080
DB_PATH        = db/data.db

# Run the application directly (outside of Docker)
run:
	@echo "Running application directly on host"
	PORT=$(PORT) DB_PATH=$(DB_PATH) go run cmd/main.go