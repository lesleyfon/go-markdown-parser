.PHONY: run dev clean deploy

# Project-specific container and image names
PROJECT_NAME=markdown-parser
DEV_IMAGE=$(PROJECT_NAME)-dev
PROD_IMAGE=$(PROJECT_NAME)
CONTAINER_NAME=markdown-api

# Default port
PORT=8080

# Production run
run:
	docker stop $(CONTAINER_NAME) 2>/dev/null || true
	docker rm $(CONTAINER_NAME) 2>/dev/null || true
	docker build -t $(PROD_IMAGE) .
	docker run --name $(CONTAINER_NAME) -p $(PORT):$(PORT) --env-file .env $(PROD_IMAGE)

# Development run with hot reload
dev:
	docker stop $(CONTAINER_NAME) 2>/dev/null || true
	docker rm $(CONTAINER_NAME) 2>/dev/null || true
	docker build -f Dockerfile.dev -t $(DEV_IMAGE) .
	docker run --name $(CONTAINER_NAME) -p $(PORT):$(PORT) --env-file .env -v $$(pwd):/app $(DEV_IMAGE)

# Clean only project-related Docker resources
clean:
	# Stop and remove project containers
	docker ps -a --filter="name=$(CONTAINER_NAME)" -q | xargs docker stop 2>/dev/null || true
	docker ps -a --filter="name=$(CONTAINER_NAME)" -q | xargs docker rm 2>/dev/null || true
	
	# Remove project images
	docker images $(PROD_IMAGE) -q | xargs docker rmi 2>/dev/null || true
	docker images $(DEV_IMAGE) -q | xargs docker rmi 2>/dev/null || true

# Check container contents
check:
	docker run -it $(PROD_IMAGE) sh -c "ls -la /app"

# Deploy to Cloud Run
deploy:
	@if [ -f deploy-config.sh ]; then \
		echo "Using local deployment configuration"; \
		source ./deploy-config.sh && \
		gcloud run deploy $(PROJECT_NAME) \
			--image gcr.io/markdown-parser-451902/$(PROJECT_NAME) \
			--platform managed \
			--allow-unauthenticated \
			--set-env-vars SECRET_KEY="$$SECRET_KEY",MONGOURI="$$MONGOURI",MONGO_DATABASE_NAME="$$MONGO_DATABASE_NAME",MONGO_SECRET="$$MONGO_SECRET",MONGO_USERNAME="$$MONGO_USERNAME" \
			--memory=2Gi \
			--timeout=10m; \
	else \
		echo "No local configuration found, deploying without env vars"; \
		gcloud run deploy $(PROJECT_NAME) \
			--image gcr.io/markdown-parser-451902/$(PROJECT_NAME) \
			--platform managed \
			--allow-unauthenticated \
			--memory=2Gi \
			--timeout=10m; \
	fi
	
