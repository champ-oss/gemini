FROM golang:1.18 AS builder

WORKDIR /app

COPY ./src/go.mod .
COPY ./src/go.sum .

RUN go mod download

COPY ./src .

RUN go build -o gemini ./cmd/


FROM debian:buster-slim

RUN apt-get update && apt-get install -y --no-install-recommends apt-utils ca-certificates

COPY --from=builder /app/gemini /

CMD ["/gemini"]