APP = k8s-controller
VERSION ?= $(shell git describe --tags --abbrev=0)
COMMIT  ?= $(shell git rev-parse --short HEAD)
DATE    ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GHCR_REGISTRY := ghcr.io/silhouetteua
ENVTEST ?= $(LOCALBIN)/setup-envtest
ENVTEST_VERSION ?= latest
LOCALBIN ?= $(shell pwd)/bin

LD_FLAGS = -X=github.com/silhouetteUA/$(APP)/cmd.Version=$(VERSION) \
           -X=github.com/silhouetteUA/$(APP)/cmd.Commit=$(COMMIT) \
           -X=github.com/silhouetteUA/$(APP)/cmd.BuildDate=$(DATE)

BUILD_FLAGS = -v -o bin/$(APP) -ldflags "$(LD_FLAGS)"

.PHONY: all build test run docker-build clean

all: build

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
ENVTEST ?= $(LOCALBIN)/setup-envtest

## Tool Versions
ENVTEST_VERSION ?= release-0.19

envtest: $(ENVTEST) ## Download setup-envtest locally if necessary.
$(ENVTEST): $(LOCALBIN)
	$(call go-install-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest,$(ENVTEST_VERSION))

format:
	gofmt -s -w ./

build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_FLAGS) main.go

test: envtest
	go install gotest.tools/gotestsum@latest
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use --bin-dir $(LOCALBIN) -p path)" gotestsum --junitfile report.xml --format testname ./... ${TEST_ARGS}

run:
	go run main.go

docker-build:
	docker build --build-arg VERSION=$(VERSION) \
	             --build-arg COMMIT=$(COMMIT) \
	             --build-arg DATE=$(DATE) \
	             -t $(GHCR_REGISTRY)/$(APP):$(VERSION)-$(COMMIT) .

clean:
	rm -rf bin
	go clean -cache