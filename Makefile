COMPOSE_PATH=docker
DEV_COMPOSE=${COMPOSE_PATH}/compose.dev.yaml

sagio.dev.start:
	@docker compose -f ${DEV_COMPOSE} up --build

sagio.dev.down:
	@docker compose -f ${DEV_COMPOSE} down

sagio.lint:
	@golangci-lint run --fix