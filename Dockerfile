FROM golang:1.20

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o ./build/youtube_production ./main.go

CMD [ "./build/youtube_production" ]


