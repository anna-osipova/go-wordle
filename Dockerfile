# backend build
FROM golang:1.17.6-alpine3.15 as go_builder
RUN apk update && apk upgrade && apk add --no-cache bash git openssh
ENV GIN_MODE=release
ENV POSTGRES_HOST=database
WORKDIR /api

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o app ./serve_words.go

# frontend build
FROM node:16.1.0-alpine AS node_builder
RUN apk add --no-cache curl
WORKDIR /client

COPY client/package-lock.json /
COPY client/package.json /
RUN npm ci
COPY client/ .
RUN npm run build

# target container
FROM alpine:latest

COPY --from=go_builder /api/app .
COPY --from=go_builder /api/words_5.txt .
COPY --from=node_builder /client/build/ ./client/build/

EXPOSE 8080
CMD ["./app"]
