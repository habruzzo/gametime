docker.build:
	docker build . -t gametime:latest --force-rm

docker.run:
	docker run -it -p 9000:9000 gametime:latest

docker.up: docker.build
	docker compose -p gametime up -d --remove-orphans

docker.dev:
	docker compose -f docker-compose-dev.yaml -p gametime-dev up -d --remove-orphans

dev.down:
	docker compose -p gametime-dev down

down:
	docker compose down

