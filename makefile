up:
	@docker compose up -d

down:
	@docker compose down	

up-test: 
	@docker compose -f docker-compose.test.yml up -d

down-test:
	@docker compose -f docker-compose.test.yml down