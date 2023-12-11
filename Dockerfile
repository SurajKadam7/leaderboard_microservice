# stage 1
FROM golang:1.20 AS build

WORKDIR /app

COPY . .

RUN go mod tidy

# use GOOS=linux GOARCH=amd64 incase of prod build on linux with architecture amd64
RUN go build -o leader_board ./main.go

# stage 2
FROM ubuntu:latest

WORKDIR /

COPY --from=build  app/leader_board /leader_board
COPY --from=build  app/config.json /config.json

ENTRYPOINT [ "/leader_board" ]
