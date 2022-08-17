docker.build:
	docker build . -t gametime:latest --force-rm

docker.clean:
	docker image prune -f

docker.run:
	docker run -it -p 9000:9000 gametime:latest

docker.up: down docker.clean docker.build
	docker compose -p gametime up -d --remove-orphans

docker.dev:
	docker compose -f docker-compose-dev.yaml -p gametime-dev up -d --remove-orphans

dev.down:
	docker compose -p gametime-dev down

down:
	docker compose -p gametime down

run:
	cd gametime; make run

dump:
	curl -v http://localhost:9000/dump -X POST -H"Authorization: "

load-backup:
	dgraph live --files ./reviews/db/data

get-dgraph:
	curl https://get.dgraph.io -sSf | bash

