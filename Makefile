LOCAL_BIN := $(CURDIR)/bin

GOOSE := $(LOCAL_BIN)/goose

DB_DSN := host=localhost port=5432 user=postgres password=postgres dbname=org_db sslmode=disable

.PHONY: install-goose
install-goose:
	@echo "Installing goose..."
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest


.PHONY: db-up
db-up:
	docker compose up -d postgres


.PHONY: wait-db
wait-db:
	@echo "Waiting for postgres..."
	until docker compose exec postgres pg_isready -U postgres; do \
		sleep 1; \
	done


.PHONY: migrate
migrate: install-goose wait-db
	@echo "Running migrations..."
	$(GOOSE) -dir ./migrations postgres "$(DB_DSN)" up


.PHONY: up
up: db-up migrate
	docker compose up --build app


.PHONY: down
down:
	docker compose down


.PHONY: restart
restart: down up