up:
	docker-compose up --build
	go run main.go

down:
	docker-compose down

install:
	go mod init
	go get ./
