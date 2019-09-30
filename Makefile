up:
	docker-compose up -d
	go run main.go

down:
	docker-compose down

install:
	go mod init
	go get ./
