
##@ Help
.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Build

GOOS ?= $(shell uname -s | tr A-Z a-z)
GOARCH ?= $(shell uname -m | sed 's/x86_64/amd64/')

build: ## Build the binary
	@echo "Building..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/$(GOOS)/$(GOARCH)/log-agent ./cmd/log-agent

##@ Docker
docker: docker-build docker-push ## Build and push

docker-build: ## Build the docker image
	@echo "Building docker image..."
	docker build -t harbor.middleware.com/middleware/log-agent:latest .

docker-push: ## Push the docker image
	@echo "Pushing docker image..."
	docker login harbor.middleware.com -u admin -p Harbor12345
	docker push harbor.middleware.com/middleware/log-agent:latest

##@ Deploy
deploy: ## Deploy the binary
	@echo "Deploying..."
	Helm install log-agent ./helm/log-agent --set image.tag=latest

undeploy: ## Undeploy the binary
	@echo "Undeploying..."
	Helm uninstall log-agent