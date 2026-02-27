SHELL := /bin/sh

COMPOSE_DEV := deployment/docker-compose.yml
COMPOSE_PROD := deployment/docker-compose-prod.yml
ENV_FILE := .env

.PHONY: dev-up dev-down dev-restart dev-logs dev-build \
	prod-up prod-down prod-restart prod-logs prod-build \
	swag swag-install

# ---------------------- Development ----------------------

dev-up:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_DEV) build --no-cache app
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_DEV) up -d

dev-down:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_DEV) down

dev-restart: dev-down dev-up


dev-logs:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_DEV) logs -f --tail=200

dev-build:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_DEV) build --pull

dev-migrate: 
	GOOSE_VERBOSE=true GOOSE_COMMAND="up" GOOSE_COMMAND_ARG="" docker compose --env-file $(ENV_FILE) -f $(COMPOSE_DEV) run --rm migrate

# ---------------------- Production ----------------------

prod-up:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_PROD) build --no-cache app
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_PROD) up -d

prod-down:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_PROD) down

prod-restart: prod-down prod-up

prod-logs:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_PROD) logs -f --tail=200

prod-build:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_PROD) build --pull

prod_migrate: 
	GOOSE_VERBOSE=true GOOSE_COMMAND="up" GOOSE_COMMAND_ARG="" docker compose --env-file $(ENV_FILE) -f $(COMPOSE_PROD) run --rm migrate