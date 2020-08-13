LDFLAGS = -ldflags "-s -w"
GOARCH = amd64
# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
BINARY_NAME=main

all: test build clean

build:
	go get -v
	$(GOBUILD) $(LDFLAGS) -o ./bin/$(BINARY_NAME) -v

sls: build
	sls deploy -v

update: build
	sls deploy --function docusign-lambda -v

deploy: sls clean

test:
	$(GOTEST) -v ./...
	
clean:
	$(GOCLEAN)
	rm -r ./bin || true

run:
	$(GOBUILD) -o $(BINARY_NAME) ${LDFLAGS} -v ./...
	./$(BINARY_NAME)