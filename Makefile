.PHONY: all snapshot-operator clean test container

ifeq ($(REGISTRY),)
	REGISTRY = quay.io/external_storage/
endif
ifeq ($(VERSION),)
	VERSION = latest
endif
IMAGE_CONTROLLER = $(REGISTRY)snapshot-operator:$(VERSION)
MUTABLE_IMAGE_CONTROLLER = $(REGISTRY)snapshot-operator:latest

all: snapshot-operator

snapshot-operator:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o _output/bin/snapshot-operator cmd/snapshot-operator/main.go

clean:
	-rm -rf _output

container:
	cp _output/snapshot-operator deploy/container/
	# Copy the root CA certificates -- cloudproviders need them
	#cp -Rf deploy/ca-certificates/* deploy/container/controller/.
	#cp -Rf deploy/ca-certificates/* deploy/container/provisioner/.
	docker build -t $(MUTABLE_IMAGE_CONTROLLER) deploy/container

test:
	go test `go list ./... | grep -v 'vendor'`
