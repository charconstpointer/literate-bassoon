FROM golang:1.14

COPY . /app
WORKDIR /app

ENV GO111MODULE=on

RUN apt-get update && apt-get install build-essential -y
RUN cd consumer && GOOS=linux go build -o cons
CMD ./consumer/cons -topic=trap -influx=http://influxdb:8086 -kafka=kafka:9092 -token=influx-golang:client