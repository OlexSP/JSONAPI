FROM golang:1.21.0-alpine3.17 AS builder

RUN mkdir /app
ADD . /app

WORKDIR /app

CMD ["go", "version"]