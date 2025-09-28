.RECIPEPREFIX := >
APP_DIR=app

.PHONY: run-api
run-api:
>cd $(APP_DIR) && API_PORT=8080 go run ./cmd/api

.PHONY: run-worker
run-worker:
>cd $(APP_DIR) && go run ./cmd/worker

.PHONY: test
test:
>cd $(APP_DIR) && go test ./...
