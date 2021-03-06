VERSION=0.0.1
USER=frperezr
SVC=microservices-demo-auth-api

BIN=$(PWD)/bin/$(SVC)
BIN_PATH=$(PWD)/bin

DB_USER=postgres
DB_NAME=postgres
DB_PASS=postgres
DSN="user=$(DB_USER) dbname=$(DB_NAME) password=$(DB_PASS) sslmode=disable"

GO ?= go
LDFLAGS='-extldflags "static" -X main.svcVersion=$(VERSION) -X main.svcName=$(SVC)'
TAGS=netgo -installsuffix netgo

migrations m:
	@echo "[migrations] Runing migrations..."
	@cd database/migrations && goose postgres $(DSN) up

clean c:
	@echo "[clean] Cleaning bin folder..."
	@rm -rf bin/

run r: migrations
	@echo "[running] Running service..."
	@go run cmd/server/main.go

build b:
	@echo "[build] Building service..."
	@cd cmd/server && $(GO) build -o $(BIN) -ldflags=$(LDFLAGS) -tags $(TAGS)

build-client bc:
	@echo "[build] Building client..."
	@cd cmd/client && $(GO) build -o $(BIN_PATH)/client -ldflags=$(LDFLAGS) -tags $(TAGS)

build-linux bl:
	@echo "[build-linux] Building service..."
	@cd cmd/server && GOOS=linux GOARCH=amd64 $(GO) build -o $(BIN) -ldflags=$(LDFLAGS) -tags $(TAGS)

build-linux-client blc:
	@echo "[build-linux-client] Building service..."
	@cd cmd/client && GOOS=linux GOARCH=amd64 $(GO) build -o $(BIN_PATH)/client -ldflags=$(LDFLAGS) -tags $(TAGS)

docker d: build-linux build-linux-client
	@echo "[build users client] Building users client..."
	@cd $(GOPATH)/src/github.com/frperezr/microservices-demo/src/users-api/cmd/client && GOOS=linux GOARCH=amd64 $(GO) build -o $(BIN_PATH)/users-api -ldflags=$(LDFLAGS) -tags $(TAGS)
	@echo "[docker] Building image..."
	@docker build -t $(USER)/$(SVC):$(VERSION) .
	@echo "[remove] Removing users client..."
	@rm bin/users-api

docker-login dl:
	@echo "[docker] Login to docker..."
	@$$(aws ecr get-login --no-include-email)

push p: docker
	@echo "[docker] pushing $(USER)/$(SVC):$(VERSION)"
	@docker tag $(USER)/$(SVC):$(VERSION) $(USER)/$(SVC):$(VERSION)
	@docker push $(USER)/$(SVC):$(VERSION)

.PHONY: migrations clean run build build-client build-linux-client docker docker-login push

