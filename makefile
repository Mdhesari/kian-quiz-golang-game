up-local:
	@docker compose -f docker-compose.local.yml watch

up:
	@docker compose up -d

down:
	@docker compose down	

prune:
	@docker image prune -f

logs:
	@docker compose logs -f -t

up-test: 
	@docker compose -f docker-compose.test.yml up -d

down-test:
	@docker compose -f docker-compose.test.yml down

reup:
	@make down && make up
restart:
	@docker compose restart