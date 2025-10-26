.PHONY: migrate migrate-up migrate-down migrate-down-steps migrate-create

migrate:
	@bash scripts/migrate $(ARGS)

migrate-up:
	@bash scripts/migrate up

migrate-down:
	@bash scripts/migrate down

migrate-down-steps:
	@bash scripts/migrate down --steps=$(number)

migrate-create:
	@bash scripts/migrate create $(name)

migrate-force:
	@bash scripts/migrate force $(version)