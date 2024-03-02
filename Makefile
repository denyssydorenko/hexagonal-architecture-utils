header:
	@echo ""
	@echo "██╗  ██╗███████╗██╗   ██╗   █████╗ ██████╗ ██╗"
	@echo "██║  ██║██╔════╝ ██╗ ██╔╝  ██╔══██╗██╔══██╗██║"
	@echo "███████║███████╗  ████╔╝   ███████║██████╔╝██║"
	@echo "██╔══██║██╔════╝ ██╔═██╗   ██╔══██║██╔═══╝ ██║"
	@echo "██║  ██║███████╗██╔╝  ██╗  ██║  ██║██║     ██║"  
	@echo "╚═╝  ╚═╝╚══════╝╚═╝   ╚═╝  ╚═╝  ╚═╝╚═╝     ╚═╝"

##@ General

help: header ## Provide help information. Example: make help
	 @awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Targets

start: ## Starts the docker-compose file provided by argument. Example: make start compose=all / Possible values to compose: all, fluent-bit, grafana, splunk
	docker-compose -f ./build/compose/$(compose)/docker-compose.yml -p service-with-$(compose) up --build --detach

stop: ## Stops the docker-compose file provided by argument. Example: make stop compose=all / Possible values to compose: all, fluent-bit, grafana, splunk
	docker-compose -f ./build/compose/$(compose)/docker-compose.yml -p service-with-$(compose) down --remove-orphans

restart: ## Restart the docker-compose file provided by argument. Example: make restart compose=all / Possible values to compose: all, fluent-bit, grafana, splunk
	docker-compose -f ./build/compose/$(compose)/docker-compose.yml -p service-with-$(compose) down --remove-orphans
	docker-compose -f ./build/compose/$(compose)/docker-compose.yml -p service-with-$(compose) up --build --detach

mocks: ## Generate mocks for unit testing. Example: make mocks
	mockgen -source=internal/ports/db.go -package db -destination=internal/adapters/db/dbmock.go

swagger:
	swag init --parseDependency --parseInternal --parseDepth 1 -g internal/adapters/http/http.go -o docs