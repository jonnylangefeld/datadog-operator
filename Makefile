VERSION ?= $(shell git describe --match 'v[0-9]*' --tags --always)
# Image URL to use all building/pushing image targets
IMG ?= jonnylangefeld/datadog-operator:$(VERSION)
CRD_OPTIONS ?= "crd:crdVersions=v1"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: manager

# Run tests
test: generate fmt vet manifests
	go test ./... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Install CRDs into a cluster
install: manifests
	kustomize build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests
	kustomize build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	cd config/manager && kustomize edit set image controller=${IMG}
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: tools
	$(GOBIN)/controller-gen $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: tools
	$(GOBIN)/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
	go generate ./...

# Build the docker image
docker-build: test
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

tools:
	@{ \
	set -e ;\
	TOOLS_TMP_DIR=$$(mktemp -d) ;\
	cp tools.go $$TOOLS_TMP_DIR ;\
	cd $$TOOLS_TMP_DIR ;\
	go mod init tmp ;\
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go get % ;\
	rm -rf $$TOOLS_TMP_DIR ;\
	}
