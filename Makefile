APP := cert-injector-webhook
IMAGE := cert-injector-webhook

.PHONY: test
test:
	@echo "\n🛠️  Running unit tests..."
	go test ./...

.PHONY: build
build:
	@echo "\n🔧  Building Go binaries..."
	GOOS=darwin GOARCH=amd64 go build -o bin/$(APP)-darwin-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o bin/$(APP)-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP)-linux-amd64 .

.PHONY: docker-build
docker-build:
	@echo "\n📦 Building $(IMAGE) Docker image..."
	docker build -t $(IMAGE):latest .

.PHONY: certs
certs:
	docker run --rm -v $(PWD)/dev:/dev debian 

# From this point `kind` is required
.PHONY: cluster
cluster:
	@echo "\n🔧 Creating Kubernetes cluster..."
	kind create cluster --config dev/manifests/kind/kind.cluster.yaml

.PHONY: delete-cluster
delete-cluster:
	@echo "\n♻️  Deleting Kubernetes cluster..."
	kind delete cluster

.PHONY: push
push: docker-build
	@echo "\n📦 Pushing $(IMAGE) into Kind's Docker daemon..."
	kind load docker-image $(IMAGE):latest

.PHONY: deploy-config
deploy-config:
	@echo "\n⚙️  Applying cluster config..."
	kubectl apply -f dev/manifests/cluster-config/

.PHONY: delete-config
delete-config:
	@echo "\n♻️  Deleting Kubernetes cluster config..."
	kubectl delete -f dev/manifests/cluster-config/

.PHONY: deploy
deploy: push delete deploy-config
	@echo "\n🚀 Deploying $(IMAGE)..."
	kubectl apply -f dev/manifests/webhook/

.PHONY: delete
delete:
	@echo "\n♻️  Deleting $(IMAGE) deployment if existing..."
	kubectl delete -f dev/manifests/webhook/ || true

.PHONY: pod
pod:
	@echo "\n🚀 Deploying test pod..."
	kubectl apply -f dev/manifests/pods/some.pod.yaml

.PHONY: delete-pod
delete-pod:
	@echo "\n♻️ Deleting test pod..."
	kubectl delete -f dev/manifests/pods/some.pod.yaml --force

.PHONY: logs
logs:
	@echo "\n🔍 Streaming $(APP) logs..."
	kubectl logs -n cert-injector -l app=$(APP) -f

.PHONY: delete-all
delete-all: delete delete-config delete-pod
