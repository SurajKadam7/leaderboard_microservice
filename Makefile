build:
	GOOS=linux GOARCH=amd64 go build -o build/youtube_production main.go

start-service: 
	docker-compose up -d 

docker-build:
	docker build -t youtube-assignment:latest .