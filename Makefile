APP = k8s-controller
VERSION ?= $(shell git describe --tags --abbrev=0)
COMMIT  ?= $(shell git rev-parse --short HEAD)
DATE    ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GHCR_REGISTRY := ghcr.io/silhouetteua

LD_FLAGS = -X=github.com/silhouetteUA/$(APP)/cmd.Version=$(VERSION) \
           -X=github.com/silhouetteUA/$(APP)/cmd.Commit=$(COMMIT) \
           -X=github.com/silhouetteUA/$(APP)/cmd.BuildDate=$(DATE)

BUILD_FLAGS = -v -o bin/$(APP) -ldflags "$(LD_FLAGS)"

.PHONY: all build test run docker-build clean

all: build

build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_FLAGS) main.go

test:
	go test ./...

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