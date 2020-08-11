  
.PHONY: all build clean run docker-test docker-build docker-push docker

GORUN = go run
GOBUILD = go build
APPNAME = authorization

all: clean build

clean:
	rm -rf bin

build: clean
	$(GOBUILD) -o bin/$(APPNAME) cmd/$(APPNAME)/*.go

docker-test:
	test $(DOCKERREPO)

docker-build: docker-test
	docker build . -t $(DOCKERREPO)

docker-push: docker-test
	docker push $(DOCKERREPO)

docker: docker-build docker-push