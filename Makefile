VPATH = .
IMAGE_NAME    ?= rpi-go-faster
IMAGE         ?= elezar/$(IMAGE_NAME)
VERSION       ?= $(shell git describe --tags --always --dirty)
TAG           ?= $(VERSION)


all: build

build: build.docker


go-faster.rpi: go-faster.go
	GOARCH=arm GOOS=linux go build -o $@ $<


.PHONY: build.docker
build.docker: go-faster.rpi Dockerfile
	docker build --rm  -t $(IMAGE):$(TAG) -f Dockerfile .

.PHONY: build.push
build.push: build.docker
	docker push "$(IMAGE):$(TAG)"

.PHONY: run
run:
	docker run --rm -ti $(IMAGE):$(TAG)