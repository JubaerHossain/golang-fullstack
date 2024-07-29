# Load environment variables from .env file
include .env
export

# Variables
APP_NAME := restaurant-saas-api
IMAGE_NAME := $(APP_NAME)
TAG := latest
DOCKER_REPO := jubaer/$(IMAGE_NAME)
KUBE_NAMESPACE := default
PORT := 3031

# Commands
rootx:
	go run cmd/rootx/main.go

install:
	go mod tidy

swag:
	swag init -g ./cmd/main.go -o docs

dev:
	go run cmd/main.go

build:
	go build -o bin/$(APP_NAME) cmd/main.go

run:
	/bin/bash -c "bin/$(APP_NAME)"

docker-build:
	docker build -t $(DOCKER_REPO):$(TAG) -f docker/Dockerfile .

docker-run:
	docker run -p $(PORT):$(PORT) $(DOCKER_REPO):$(TAG)

docker-compose-up:
	docker-compose -f docker-compose.yaml up -d

docker-compose-down:
	docker-compose -f docker-compose.yaml down

docker-prune:
	docker system prune -f

# Kubernetes related commands
deploy:
	kubectl apply -f k8s/configmap.yaml --validate=false
	kubectl apply -f k8s/secret.yaml --validate=false
	kubectl apply -f k8s/postgres-pv.yaml --validate=false
	kubectl apply -f k8s/postgres-pvc.yaml --validate=false
	kubectl apply -f k8s/postgres-deployment.yaml --validate=false
	kubectl apply -f k8s/postgres-service.yaml --validate=false
	kubectl apply -f k8s/redis-pv.yaml --validate=false
	kubectl apply -f k8s/redis-pvc.yaml --validate=false
	kubectl apply -f k8s/redis-deployment.yaml --validate=false
	kubectl apply -f k8s/redis-service.yaml --validate=false
	kubectl apply -f k8s/deployment.yaml --validate=false
	kubectl apply -f k8s/service.yaml --validate=false
	kubectl apply -f k8s/ingress.yaml --validate=false
	kubectl apply -f k8s/horizontal-pod-autoscaler.yaml --validate=false





re-deploy: deploy

cleanup:
	kubectl delete -f k8s/configmap.yaml || true
	kubectl delete -f k8s/secret.yaml || true
	kubectl delete -f k8s/deployment.yaml || true
	kubectl delete -f k8s/service.yaml || true
	kubectl delete -f k8s/ingress.yaml || true
	kubectl delete -f k8s/horizontal-pod-autoscaler.yaml || true

# Build and push Docker image
build-push:
	docker build -t $(DOCKER_REPO):$(TAG) -f docker/Dockerfile .
	docker push $(DOCKER_REPO):$(TAG)

# Build, push Docker image, and deploy
release: build-push deploy

# Additional utilities
proto:
	@./scripts/proto.sh domain/foods/infrastructure/transport/grpc/proto

.PHONY: rootx install swag dev build run docker-build docker-run docker-compose-up docker-compose-down docker-prune deploy re-deploy cleanup build-push release proto
