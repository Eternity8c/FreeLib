FROM golang:alpine3.23 AS dependencies
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 

FROM dependencies AS builder
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:3.23
RUN apk add --no-cache wget
WORKDIR /app
COPY --from=builder /app/server /app/server

CMD [ "/app/server" ]