build:
	docker build --no-cache --tag chaseisabelle/pipeline-example:latest .

run:
	docker run --rm -it chaseisabelle/pipeline-example:latest

up:
	make build
	make run

rmi:
	docker rmi -f chaseisabelle/pipeline-example:latest
