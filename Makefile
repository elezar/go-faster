VPATH = .
IMAGE_NAME    ?= go-faster
RPI_IMAGE     ?= elezar/rpi-$(IMAGE_NAME)
IMAGE         ?= elezar/$(IMAGE_NAME)
VERSION       ?= $(shell git describe --tags --always --dirty)
TAG           ?= $(VERSION)


all: build

build: build.docker

build.docker: build.docker.rpi build.docker.amd64


go-faster.rpi: go-faster.go
	GOARCH=arm GOOS=linux go build -o $@ $<

go-faster.amd64: go-faster.go
	GOARCH=amd64 GOOS=linux go build -o $@ $<

.PHONY: build.docker.rpi
build.docker.rpi: go-faster.rpi Dockerfile.rpi
	docker build --rm  -t $(RPI_IMAGE):$(TAG) -f Dockerfile.rpi .

.PHONY: build.docker.amd64
build.docker.amd64: go-faster.amd64 Dockerfile.amd64
	docker build --rm  -t $(IMAGE):$(TAG) -f Dockerfile.amd64 .

.PHONY: tag.latest
tag.latest: build.docker
	docker tag $(IMAGE):$(TAG) $(IMAGE):latest
	docker tag $(RPI_IMAGE):$(TAG) $(RPI_IMAGE):latest

.PHONY: build.push
build.push: build.docker tag.latest
	docker push "$(IMAGE):$(TAG)"
	docker push "$(IMAGE):latest"
	docker push "$(RPI_IMAGE):$(TAG)"
	docker push "$(RPI_IMAGE):latest"

.PHONY: run.rpi
run.rpi: build.docker.rpi
	docker run --rm -ti $(RPI_IMAGE):$(TAG)

.PHONY: run
run: build.docker.amd64
	docker run --rm -ti $(IMAGE):$(TAG)


clean:
	docker rmi -f $(RPI_IMAGE):$(TAG)
	docker rmi -f $(IMAGE):$(TAG)
	rm -f go-faster.rpi go-faster.amd64
