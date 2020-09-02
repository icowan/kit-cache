APPNAME = kit-cache
BIN = $(GOPATH)/bin
GOCMD = /usr/local/go/bin/go
GOBUILD = $(GOCMD) build
GOINSTALL = $(GOCMD) install
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GORUN = $(GOCMD) run
BINARY_UNIX = $(BIN)/$(APPNAME)
PID = .pid
HUB_ADDR = operations-virtual-local.repo.yrd.creditease.corp
DOCKER_USER =
DOCKER_PWD =
TAG = v0.0.01-test
NAMESPACE = operations
PWD = $(shell pwd)

start:
	$(BIN)/$(APPNAME) start -p :8080 -c /etc/yxlive/app.cfg -i remote & echo $$! > $(PID)

restart:
	@echo restart the app...
	@kill `cat $(PID)` || true
	$(BIN)/$(APPNAME) start -p :8080 -c /etc/yxlive/app.cfg -i remote & echo $$! > $(PID)

install:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOINSTALL) -v

stop:
	@kill `cat $(PID)` || true

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

login:
	docker login -u $(DOCKER_USER) -p $(DOCKER_PWD) $(HUB_ADDR)

docker-build:
	docker build --rm -t $(APPNAME):$(TAG) .

docker-run:
	docker run -it --rm -p 8080:8080 -v $(PWD)/app.cfg:/etc/yxlive/app.cfg $(APPNAME):$(TAG)

push:
	docker image tag $(APPNAME):$(TAG) $(HUB_ADDR)/$(NAMESPACE)/$(APPNAME):$(TAG)
	docker push $(HUB_ADDR)/$(NAMESPACE)/$(APPNAME):$(TAG)

run:
	#GO111MODULE=on $(GOGET) -v -insecure gitlab.creditease.corp/yxlive/types
	#GO111MODULE=on $(GOGET) -v -insecure gitlab.creditease.corp/pkg/api/alert
	GO111MODULE=on $(GORUN) ./cmd/main.go start -p :8080 -g :8082 -c ./app.dev.cfg