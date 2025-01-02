# Docker-related variables
DOCKERFILE     = Dockerfile
IMAGE_NAME     = forum-image
CONTAINER_NAME = forum-container

# Application-related variables
APP_DIR        = $(PWD)
PORT           = 8080
DB_PATH        = db/data.db

# Build the Docker image
build: 
	@echo "Building Docker image: $(IMAGE_NAME)"
	docker build -f $(DOCKERFILE) -t $(IMAGE_NAME) .

# Run the Docker container
run:
	@echo "Running Docker container: $(CONTAINER_NAME)"
	@echo $(PORT)
	 docker run --name $(CONTAINER_NAME) -p $(PORT):$(PORT) $(IMAGE_NAME)

# Stop and remove the Docker container
stop:
	@echo "Stopping and removing Docker container: $(CONTAINER_NAME)"
	 docker stop $(CONTAINER_NAME) || true

# Clean up files and Docker resources
clean:
	@echo "Cleaning up temporary files and Docker resources"
	rm -rf forum || true
	docker rm -f $(CONTAINER_NAME) || true
	docker rmi -f $(IMAGE_NAME) || true

# Push changes to Git repository
push: clean
	@echo "Committing and pushing changes to Git"
	@read -p "Enter commit message: " msg; \
	git add .; \
	git commit -m "$$msg"; \
	git push

# Combined target for stopping, cleaning, building, and running
all: stop clean build run

# Run the application directly (outside of Docker)
run-app:
	@echo "Running application directly on host"
	PORT=$(PORT) DB_PATH=$(DB_PATH) go run cmd/main.go