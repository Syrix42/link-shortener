Link Shortener

A Go + Fiber API with Docker Compose dev/prod workflows and Swagger docs via `swaggo/swag`.

Make Targets

- `make dev-up`: Start dev stack (`deployment/docker-compose.yml`)
- `make dev-down`: Stop dev stack
- `make dev-restart`: Restart dev stack
- `make dev-logs`: Tail dev logs
- `make dev-build`: Build dev images
- `make dev-migrate`: Migrate Database
- `make prod-up`: Start prod stack (`deployment/docker-compose-prod.yml`)
- `make prod-down`: Stop prod stack
- `make prod-restart`: Restart prod stack
- `make prod-logs`: Tail prod logs
- `make prod-build`: Build prod images
- `make prod-migrate`: Migrate Database

Swagger

- Open in browser: start the app, then visit `/swagger/index.html`
