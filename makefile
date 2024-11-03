up-local:
	@docker compose -f docker-compose.local.yml watch -d

up:
	@docker compose up -d

down:
	@docker compose down	

prune:
	@docker image prune -f

log:
	@docker compose logs -f -t

up-test: 
	@docker compose -f docker-compose.test.yml up -d

down-test:
	@docker compose -f docker-compose.test.yml down