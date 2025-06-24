APP = k8s-controller
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT  ?= $(shell git rev-parse --short HEAD)
DATE    ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

LD_FLAGS = -X=github.com/silhouetteUA/$(APP)/cmd.Version=$(VERSION) \
           -X=github.com/silhouetteUA/$(APP)/cmd.Commit=$(COMMIT) \
           -X=github.com/silhouetteUA/$(APP)/cmd.BuildDate=$(DATE)

BUILD_FLAGS = -v -o $(APP) -ldflags "$(LD_FLAGS)"

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
	             -t $(APP):latest .

clean:
	rm -f $(APP)
	go clean -cache