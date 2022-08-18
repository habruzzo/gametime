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
	docker run -it --name dgraph-live --rm -v "`pwd`/reviews/db/data:/dump" --net gametime-dev_default dgraph/dgraph:latest dgraph live -f /dump -a alpha:9080 -z zero:5080

box.load-backup:
	docker run -it --name dgraph-live --rm -v "`pwd`/reviews/db/data:/dump" --net gametime-box_default dgraph/dgraph:latest dgraph live -f /dump -a alpha:9080 -z zero:5080

box.docker.up: box.down docker.clean docker.build
	docker-compose -p gametime up -d --remove-orphans

box.docker.dev:
	docker-compose -f docker-compose-box.yaml -p gametime-box up -d --remove-orphans

box.dev.down:
	docker-compose -p gametime-box down

box.down:
	docker-compose -p gametime down

box.backup: dump
	cd gametime/reviews
	git add .
	git commit -m "backup `date`"
	git push

build.review:
	./tools/review_builder.py

push.review:
	cd ~/gametime/reviews
	git add .
	git commit -m "add review"
	git push
	cd ..
	git submodule update


review: build.review push.review

prep-start:
	./tools/prep-box.sh start

prep-restart:
	./tools/prep-box.sh restart

tunnel:
	ssh
