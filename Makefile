up:
	docker-compose up --build

down:
	docker-compose down

install:
	go mod init
	go get ./
