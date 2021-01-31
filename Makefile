PROJ_NAME = azlyrics-fetcher
generate:
	oapi-codegen -generate "chi-server" --package="app" openapi.yaml  > ./internal/app/server.gen.go
	oapi-codegen -generate "types" --package="azlyrics" openapi.yaml  > ./models.gen.go

restart: stop start

start:
	docker-compose up -d

debug:
	gebug start

stop:
	docker-compose down

build:
	docker-compose build --parallel

logs:
	docker-compose logs -f

test:
	go test -v ./...

check:
	docker run \
		-v $(shell pwd):/go/src/github.com/johnmackenzie91/azlyrics-fetcher \
		golangci/golangci-lint:v1.36.0 \
		golangci-lint -v run ./src/github.com/johnmackenzie91/azlyrics-fetcher/...