up:
	docker-compose up -d

down:
	docker-compose down --remove-orphans --volumes

build:
	docker build . -t localhost:5000/jean/rinha-backend:latest && docker push localhost:5000/jean/rinha-backend:latest

stress:
	./stress-test/run.sh

test:
	go test -v ./...

up-dev:
	docker-compose -f docker-compose-dev.yml up -d

down-dev:
	docker-compose -f docker-compose-dev.yml down --remove-orphans --volumes
