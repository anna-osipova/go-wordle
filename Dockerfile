FROM golang:1.17.6-alpine3.15 as builder
RUN apk update && apk upgrade && apk add --no-cache bash git openssh
ENV GIN_MODE=release
ENV POSTGRES_HOST=database
WORKDIR /api

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o app ./serve_words.go

FROM alpine:latest

COPY --from=builder /api/app .
COPY --from=builder /api/words_5.txt .

EXPOSE 8080
CMD ["./app"]
