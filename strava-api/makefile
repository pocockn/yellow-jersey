BINARY_UNIX=$(REPONAME)_unix
PACKAGES ?= "./..."
DOCKERNAME = "pococknick91"
REPONAME ?= "yellow-jersey"
IMG ?= ${DOCKERNAME}/${REPONAME}:${VERSION}
LATEST ?= ${DOCKERNAME}/${REPONAME}:latest
# Use tag for current commit, or the most recent tag and short commit SHA
VERSION="0.1.0"
# $(shell VERSION=$$(git tag --points-at HEAD);SHA_VERSION=$$(git describe --tags --abbrev=0 $$(git rev-list --tags --max-count=1))-$$(git rev-parse --short HEAD);echo $${VERSION:=$${SHA_VERSION}})
DEFAULT: test

build:
	@GO111MODULE=on go build "${PACKAGES}"

build-image:
	@docker build -t ${IMG} .

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_UNIX) -v

## clean-mock: removes all generated mocks
clean-mock:
	@find . -iname '*_mock.go' -exec rm {} \;

## generate: runs go generate
generate:
	go install github.com/golang/mock/mockgen@v1.6.0
	go generate -v ./...

install:
	@echo "=> Install dependencies"
	@GO111MODULE=on go mod download

push-to-registry:
	@docker login -u ${DOCKER_USER} -p ${DOCKER_PASS}
	@docker build -t ${IMG} .
	@docker tag ${IMG} ${LATEST}
	echo "=> Pushing ${IMG} & ${LATEST} to docker"
	@docker push ${DOCKERNAME}/${REPONAME}

run:
	@go build -o strava-api -ldflags="-X 'main.Version=${VERSION}'" ./cmd/main.go
	@ENV=dev AWS_REGION=eu-west-1 ./strava-api

test:
	@GO111MODULE=on go test "${PACKAGES}" -cover

regenerate: clean-mock generate

vet:
	@@GO111MODULE=on go vet "${PACKAGES}"

fix:
	@go fix ./...

lint:
	golangci-lint --version
	golangci-lint cache clean
	golangci-lint run --timeout=5m