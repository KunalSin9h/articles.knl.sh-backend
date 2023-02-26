run:
	DB=./db/dev.db \
	PORT=5000 \
	PASSWORD=1 \
	go run ./src/main.go

build:
	docker build -t ghcr.io/kunalsin9h/article-back:latest .